package biz

import (
	"bytes"
	"ehc-hsm-server/utils"
	"encoding/binary"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
)

//验证健康卡或解密时间
type S4Req struct {
	ReqCode           string
	Var1              int32
	Var2              int32
	Var3Ver           string //String.format("K%04d", 102 + ver * 2);
	Var4IndexidLength int32  // 主索引.length() / 32
	Var5IndexidFactor string
	Var6              int32  //固定值0
	Var7              string //空值
	Var8PaddingMode   int32  //PaddingMode
	Var9DataLength    int32  //（证件类型+证件号码）长度
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

func (c *S4) Decode(buf []byte) error {
	in := bytes.NewBuffer(buf)
	var err error = nil
	c.req.ReqCode, _ = utils.ReadString(in, 2)
	c.req.Var1, err = utils.ReadInt(in, 2)
	c.req.Var2, err = utils.ReadIntHex(in, 3)
	c.req.Var3Ver, err = utils.ReadString(in, 5)
	c.req.Var4IndexidLength, err = utils.ReadInt(in, 2)
	c.req.Var5IndexidFactor, err = utils.ReadString(in, c.req.Var4IndexidLength*32)
	c.req.Var6, err = utils.ReadInt(in, 2)
	//req.Var7 = utils.ReadString(in, 3);
	c.req.Var8PaddingMode, err = utils.ReadInt(in, 2) //ANSI_X919("ANSIX919PADDING", 2), 填充模式
	c.req.Var9DataLength, err = utils.ReadIntHex(in, 4)
	c.req.Var9Data, err = utils.ReadBytes(in, c.req.Var9DataLength)
	return err
}

func (c *S4) setData(data []byte) {
	if data != nil {
		c.res.data = data
		c.res.dataLength = len(data)
	} else {
		c.res.data = nil
		c.res.dataLength = 0
	}
}

func (c *S4) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write([]byte(c.res.ReqCode))
	buf.Write([]byte(c.res.ErrCode))

	//String hex = Integer.toHexString(DataLength);
	//final byte[] len = Strings.padStart(hex, 4, '0').getBytes(Charsets.ISO_8859_1);
	var bufLen = make([]byte, 2)
	binary.BigEndian.PutUint16(bufLen, uint16(c.res.dataLength))
	lenStr := hex.EncodeToString(bufLen)

	buf.Write([]byte(lenStr))
	if c.res.data != nil {
		buf.Write(c.res.data)
	}
	return buf.Bytes(), nil
}

func (c *S4) Handle() error {
	c.res.ErrCode = "00"

	if c.req.Var4IndexidLength > 0 {
		//验证健康卡
		indexidFactor := c.req.Var5IndexidFactor
		data, _ := hex.DecodeString(indexidFactor)
		safeKey, _ := hex.DecodeString(SAFE_KEY)
		ehcIdKeys, _ := utils.Sm4EncryptNoPadding(safeKey, data)
		inputData := c.req.Var9Data
		log.Printf("用户身份认证密钥=%s", hex.EncodeToString(ehcIdKeys))
		//如果是nopadding则传入的字节数据长度应该是16的倍数
		ehcIds, err := utils.Sm4DecryptNoPadding(ehcIdKeys, inputData)
		if err != nil {
			return err
		}
		c.setData(ehcIds)
		log.Printf("验证健康卡=%s", hex.EncodeToString(ehcIds))
	} else if c.req.Var4IndexidLength == 0 {
		//// 解密时间
		facKeys, _ := hex.DecodeString(PROTECT_KEY)
		var var9Data []byte = c.req.Var9Data
		log.Println("解密时间in={}", hex.EncodeToString(var9Data))
		outData, err := utils.Sm4DecryptNoPadding(facKeys, var9Data)
		if err != nil {
			return err
		}
		c.setData(outData)
		log.Println("解密时间out={}", hex.EncodeToString(outData))
	}
	return nil
}
