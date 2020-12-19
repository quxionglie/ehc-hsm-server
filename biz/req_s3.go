package biz

import (
	"bytes"
	"ehc-hsm-server/pkg"
	"ehc-hsm-server/pkg/rwbytes"
	"encoding/binary"
	hex "encoding/hex"
	log "github.com/sirupsen/logrus"
	"strings"
)

//生成健康卡或加密时间
// "%02X%03X%s%02X%s%02X%s%02X%04X"
// var1 %02X  0
// var2 %03X  10
// var3 %s    版本 ver K%04d
// var4 %02X
// var5 %s
// var6 %02X
// var7 %s
// var8 %02X
// var9 %04X
type S3Req struct {
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

type S3Res struct {
	ReqCode    string //2A
	ErrCode    string //2A
	dataLength int    //2N
	data       []byte
}

type S3 struct {
	req *S3Req
	res *S3Res
}

func NewS3() *S3 {
	req := S3Req{}
	res := S3Res{ReqCode: "S4", ErrCode: "00"}
	return &S3{&req, &res}
}

func (req *S3Req) Decode(buf []byte) error {
	in := bytes.NewBuffer(buf)
	log.Printf("请求内容S3Req=%s", strings.ToUpper(hex.EncodeToString(in.Bytes())))
	var err error
	req.ReqCode, _ = rwbytes.ReadString(in, 2)
	req.Var1, _ = rwbytes.ReadInt(in, 2)
	req.Var2, _ = rwbytes.ReadIntHex(in, 3)
	req.Var3Ver, _ = rwbytes.ReadString(in, 5)
	req.Var4IndexidLength, _ = rwbytes.ReadInt(in, 2)
	req.Var5IndexidFactor, _ = rwbytes.ReadString(in, req.Var4IndexidLength*32)
	req.Var6, _ = rwbytes.ReadInt(in, 2)
	//Req.Var7 = pkg.ReadString(in, 3);
	req.Var8PaddingMode, _ = rwbytes.ReadInt(in, 2) //ANSI_X919("ANSIX919PADDING", 2), 填充模式
	req.Var9DataLength, _ = rwbytes.ReadIntHex(in, 4)
	req.Var9Data, err = rwbytes.ReadBytes(in, req.Var9DataLength)
	if err != nil {
		log.Printf("解析出错=%v", err)
	}
	return err
}

func (req *S3Req) Encode(buf []byte) ([]byte, error) {
	in := new(bytes.Buffer)
	rwbytes.WriteString(in, 2, req.ReqCode)
	rwbytes.WriteInt(in, 2, req.Var1)
	rwbytes.WriteIntHex(in, 3, req.Var2)
	rwbytes.WriteString(in, 5, req.Var3Ver)
	rwbytes.WriteInt(in, 2, req.Var4IndexidLength)
	rwbytes.WriteString(in, req.Var4IndexidLength*32, req.Var5IndexidFactor)
	rwbytes.WriteInt(in, 2, req.Var6)
	//Req.Var7 = pkg.WriteString(in, 3);
	rwbytes.WriteInt(in, 2, req.Var8PaddingMode) //ANSI_X919("ANSIX919PADDING", 2), 填充模式
	rwbytes.WriteIntHex(in, 4, req.Var9DataLength)
	rwbytes.WriteBytes(in, req.Var9DataLength, req.Var9Data)
	return in.Bytes(), nil
}

func (res *S3Res) setData(data []byte) {
	if data != nil {
		res.data = data
		res.dataLength = len(data)
	} else {
		res.data = nil
		res.dataLength = 0
	}
}

func (res *S3Res) Decode(b []byte) error {
	in := bytes.NewBuffer(b)
	res.ReqCode, _ = rwbytes.ReadString(in, 2)
	res.ErrCode, _ = rwbytes.ReadString(in, 2)
	if "00" == res.ErrCode {
		var bufLen = make([]byte, 2)
		in.Read(bufLen)
		binary.BigEndian.Uint16(bufLen)
		res.dataLength = int(binary.BigEndian.Uint16(bufLen))
		res.data, _ = rwbytes.ReadBytes(in, res.dataLength)
	}
	return nil
}

func (res *S3Res) Encode() ([]byte, error) {
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

func (c *S3) Decode(buf []byte) error {
	return c.req.Decode(buf)
}
func (c *S3) Encode() ([]byte, error) {
	return c.res.Encode()
}

func (c *S3) Handle() error {
	if c.req.Var4IndexidLength > 0 {
		// 生成健康卡
		// name=张秦&idNo=440303199405109011
		// 主索引：1E5E48D529BC88CB134E06B8C18D67578372F51D490D4A0230298163E94CB017
		// 健康卡：1D341DF86A1B41443A6531797659C2A05B856727FEE452F6BD5E89F2D41B74CB
		indexidFactor := c.req.Var5IndexidFactor
		data, _ := hex.DecodeString(indexidFactor)
		safeKey, _ := hex.DecodeString(SAFE_KEY)
		ehcIdKeys, _ := pkg.Sm4EncryptNoPadding(safeKey, data)
		inputIdNo := c.req.Var9Data
		log.Printf("用户身份认证密钥=%s", hex.EncodeToString(ehcIdKeys))
		//如果是nopadding则传入的字节数据长度应该是16的倍数
		idNoBytes := make([]byte, 32)
		for i := 0; i < len(idNoBytes); i++ {
			if i < len(inputIdNo) {
				idNoBytes[i] = inputIdNo[i]
			}
		}
		ehcIds, err1 := pkg.Sm4EncryptNoPadding(ehcIdKeys, idNoBytes)
		if err1 != nil {
			log.Println("出错=", err1)
		}
		c.res.setData(ehcIds)
		log.Printf("健康卡id=%s", strings.ToUpper(hex.EncodeToString(ehcIds)))
	} else if c.req.Var4IndexidLength == 0 {
		//// 加密时间
		facKeys, _ := hex.DecodeString(PROTECT_KEY)
		var var9Data []byte = c.req.Var9Data
		log.Println("加密时间in={}", strings.ToUpper(hex.EncodeToString(var9Data)))
		//如果是nopadding则传入的字节数据长度应该是16的倍数
		var tmpBytes []byte = make([]byte, 16)
		for i := 0; i < len(tmpBytes); i++ {
			if i < len(var9Data) {
				tmpBytes[i] = var9Data[i]
			}
		}
		outData, _ := pkg.Sm4EncryptNoPadding(facKeys, tmpBytes)
		c.res.setData(outData)
		log.Println("加密时间out={}", hex.EncodeToString(outData))
	}
	return nil

}
