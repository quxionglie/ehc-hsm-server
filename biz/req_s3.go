package biz

import (
	"bytes"
	"ehc-hsm-server/utils"
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

func (c *S3) Decode(buf []byte) error {
	in := bytes.NewBuffer(buf)
	log.Printf("请求内容S3Req=%s", strings.ToUpper(hex.EncodeToString(in.Bytes())))
	var err error
	c.req.ReqCode, _ = utils.ReadString(in, 2)
	c.req.Var1, _ = utils.ReadInt(in, 2)
	c.req.Var2, _ = utils.ReadIntHex(in, 3)
	c.req.Var3Ver, _ = utils.ReadString(in, 5)
	c.req.Var4IndexidLength, _ = utils.ReadInt(in, 2)
	c.req.Var5IndexidFactor, _ = utils.ReadString(in, c.req.Var4IndexidLength*32)
	c.req.Var6, _ = utils.ReadInt(in, 2)
	//req.Var7 = utils.ReadString(in, 3);
	c.req.Var8PaddingMode, _ = utils.ReadInt(in, 2) //ANSI_X919("ANSIX919PADDING", 2), 填充模式
	c.req.Var9DataLength, _ = utils.ReadIntHex(in, 4)
	c.req.Var9Data, err = utils.ReadBytes(in, c.req.Var9DataLength)
	if err != nil {
		log.Printf("解析出错=%v", err)
	}
	return err
}

func (c *S3) setData(data []byte) {
	if data != nil {
		c.res.data = data
		c.res.dataLength = len(data)
	} else {
		c.res.data = nil
		c.res.dataLength = 0
	}
}

func (c *S3) Encode() ([]byte, error) {
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

func (c *S3) Handle() error {
	if c.req.Var4IndexidLength > 0 {
		// 生成健康卡
		// name=张秦&idNo=440303199405109011
		// 主索引：1E5E48D529BC88CB134E06B8C18D67578372F51D490D4A0230298163E94CB017
		// 健康卡：1D341DF86A1B41443A6531797659C2A05B856727FEE452F6BD5E89F2D41B74CB
		indexidFactor := c.req.Var5IndexidFactor
		data, _ := hex.DecodeString(indexidFactor)
		safeKey, _ := hex.DecodeString(SAFE_KEY)
		ehcIdKeys, _ := utils.Sm4EncryptNoPadding(safeKey, data)
		inputIdNo := c.req.Var9Data
		log.Printf("用户身份认证密钥=%s", hex.EncodeToString(ehcIdKeys))
		//如果是nopadding则传入的字节数据长度应该是16的倍数
		idNoBytes := make([]byte, 32)
		for i := 0; i < len(idNoBytes); i++ {
			if i < len(inputIdNo) {
				idNoBytes[i] = inputIdNo[i]
			}
		}
		ehcIds, err1 := utils.Sm4EncryptNoPadding(ehcIdKeys, idNoBytes)
		if err1 != nil {
			log.Println("出错=", err1)
		}
		c.setData(ehcIds)
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
		outData, _ := utils.Sm4EncryptNoPadding(facKeys, tmpBytes)
		c.setData(outData)
		log.Println("加密时间out={}", hex.EncodeToString(outData))
	}
	return nil

}
