package biz

import (
	"bytes"
	"ehc-hsm-server/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tjfoc/gmsm/sm3"
)

//3c生成摘要
type Digest3cReq struct {
	ReqCode          string
	AlgorithmPattern string
	DataLength       int32
	Data             []byte
	EndStr           string
	Var3Length       int32
	Var4Length       int32
}

//3c生成摘要
type Digest3cRes struct {
	ReqCode    string //2A
	ErrCode    string //2A
	DataLength int    //2N
	Data       []byte
}

type Digest3c struct {
	req *Digest3cReq
	res *Digest3cRes
}

//
//func NewDigest3cHandle() BizHandle{
//
//}

func NewDigest3c() *Digest3c {
	req := Digest3cReq{}
	res := Digest3cRes{ReqCode: "3D", ErrCode: "00"}
	return &Digest3c{&req, &res}
}

func (c *Digest3c) Decode(buf []byte) error {
	in := bytes.NewBuffer(buf)

	var err error = nil
	c.req.ReqCode, err = utils.ReadString(in, 2)
	c.req.AlgorithmPattern, err = utils.ReadString(in, 2)
	c.req.DataLength, err = utils.ReadInt(in, 4)
	c.req.Data, err = utils.ReadBytes(in, c.req.DataLength)
	c.req.EndStr, err = utils.ReadString(in, 1)
	c.req.Var3Length, err = utils.ReadInt(in, 4)
	c.req.Var4Length, err = utils.ReadInt(in, 4)
	log.Printf("请求参数：%v", c.req)
	return err
}

func (c *Digest3c) setData(data []byte) {
	if data != nil {
		c.res.Data = data
		c.res.DataLength = len(data)
	} else {
		c.res.Data = nil
		c.res.DataLength = 0
	}
}

func (c *Digest3c) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write([]byte(c.res.ReqCode))
	buf.Write([]byte(c.res.ErrCode))

	//byte[] len = Strings.padStart(DataLength + "", 2, '0').getBytes(Charsets.ISO_8859_1);
	lenStr := fmt.Sprintf("%02d", c.res.DataLength)
	buf.Write([]byte(lenStr))
	if c.res.Data != nil {
		buf.Write(c.res.Data)
	}
	return buf.Bytes(), nil
}

func (c *Digest3c) Handle() error {
	c.res.ErrCode = "00"
	h := sm3.New()
	h.Write(c.req.Data)
	sum := h.Sum(nil)
	c.setData(sum)
	return nil
}
