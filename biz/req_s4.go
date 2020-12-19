package biz

import (
	"bytes"
	"ehc-hsm-server/pkg"
	"ehc-hsm-server/pkg/rwbytes"
	"encoding/binary"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
)

//验证健康卡或解密时间
type S4Req struct {
	ReqCode           string
	Var1              int
	Var2              int
	Var3Ver           string //String.format("K%04d", 102 + ver * 2);
	Var4IndexidLength int    // 主索引.length() / 32
	Var5IndexidFactor string
	Var6              int    //固定值0
	Var7              string //空值
	Var8PaddingMode   int    //PaddingMode
	Var9DataLength    int    //（证件类型+证件号码）长度
	Var9Data          []byte //（证件类型+证件号码）
	Var10             string //空值
}

type S4Res struct {
	ReqCode    string //2A
	ErrCode    string //2A
	dataLength int    //2N
	data       []byte
}

type S4 struct {
	req *S4Req
	res *S4Res
}

func NewS4() *S4 {
	req := S4Req{}
	res := S4Res{ReqCode: "S5", ErrCode: "00"}
	return &S4{&req, &res}
}

func (req *S4Req) Decode(buf []byte) error {
	in := bytes.NewBuffer(buf)
	var err error = nil
	req.ReqCode, _ = rwbytes.ReadString(in, 2)
	req.Var1, err = rwbytes.ReadInt(in, 2)
	req.Var2, err = rwbytes.ReadIntHex(in, 3)
	req.Var3Ver, err = rwbytes.ReadString(in, 5)
	req.Var4IndexidLength, err = rwbytes.ReadInt(in, 2)
	req.Var5IndexidFactor, err = rwbytes.ReadString(in, req.Var4IndexidLength*32)
	req.Var6, err = rwbytes.ReadInt(in, 2)
	//Req.Var7 = pkg.ReadString(in, 3);
	req.Var8PaddingMode, err = rwbytes.ReadInt(in, 2) //ANSI_X919("ANSIX919PADDING", 2), 填充模式
	req.Var9DataLength, err = rwbytes.ReadIntHex(in, 4)
	req.Var9Data, err = rwbytes.ReadBytes(in, req.Var9DataLength)
	return err
}

func (res *S4Res) setData(data []byte) {
	if data != nil {
		res.data = data
		res.dataLength = len(data)
	} else {
		res.data = nil
		res.dataLength = 0
	}
}

func (res *S4Res) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write([]byte(res.ReqCode))
	buf.Write([]byte(res.ErrCode))

	//String hex = Integer.toHexString(DataLength);
	//final byte[] len = Strings.padStart(hex, 4, '0').getBytes(Charsets.ISO_8859_1);
	var bufLen = make([]byte, 2)
	binary.BigEndian.PutUint16(bufLen, uint16(res.dataLength))
	lenStr := hex.EncodeToString(bufLen)

	buf.Write([]byte(lenStr))
	if res.data != nil {
		buf.Write(res.data)
	}
	return buf.Bytes(), nil
}

func (c *S4) Decode(buf []byte) error {
	return c.req.Decode(buf)
}
func (c *S4) Encode() ([]byte, error) {
	return c.res.Encode()
}

func (c *S4) Handle() error {
	c.res.ErrCode = "00"

	if c.req.Var4IndexidLength > 0 {
		//验证健康卡
		indexidFactor := c.req.Var5IndexidFactor
		data, _ := hex.DecodeString(indexidFactor)
		safeKey, _ := hex.DecodeString(SAFE_KEY)
		ehcIdKeys, _ := pkg.Sm4EncryptNoPadding(safeKey, data)
		inputData := c.req.Var9Data
		log.Printf("用户身份认证密钥=%s", hex.EncodeToString(ehcIdKeys))
		//如果是nopadding则传入的字节数据长度应该是16的倍数
		ehcIds, err := pkg.Sm4DecryptNoPadding(ehcIdKeys, inputData)
		if err != nil {
			return err
		}
		c.res.setData(ehcIds)
		log.Printf("验证健康卡=%s", hex.EncodeToString(ehcIds))
	} else if c.req.Var4IndexidLength == 0 {
		//// 解密时间
		facKeys, _ := hex.DecodeString(PROTECT_KEY)
		var var9Data []byte = c.req.Var9Data
		log.Println("解密时间in={}", hex.EncodeToString(var9Data))
		outData, err := pkg.Sm4DecryptNoPadding(facKeys, var9Data)
		if err != nil {
			return err
		}
		c.res.setData(outData)
		log.Println("解密时间out={}", hex.EncodeToString(outData))
	}
	return nil
}
