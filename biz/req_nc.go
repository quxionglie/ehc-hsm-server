package biz

import (
	"bytes"
	"ehc-hsm-server/utils"
)

type NcReq struct {
	ReqCode string
}

type NcRes struct {
	ReqCode string
	ErrCode string
	Data    string
}

type Nc struct {
	req *NcReq
	res *NcRes
}

func NewNc() *Nc {
	req := NcReq{}
	res := NcRes{ReqCode: "ND", ErrCode: "00"}
	return &Nc{&req, &res}
}

func (c *Nc) Decode(in []byte) (err error) {
	buf := bytes.NewBuffer(in)
	c.req.ReqCode, err = utils.ReadString(buf, 2)
	return
}

func (c *Nc) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write([]byte(c.res.ReqCode))
	buf.Write([]byte(c.res.ErrCode))
	if c.res.Data != "" {
		buf.Write([]byte(c.res.Data))
	}
	return buf.Bytes(), nil
}

func (c *Nc) Handle() error {
	//req := NewNcReq()
	//req.Decode(buf)
	//log.Printf("请求参数：%v", req)
	c.res.ErrCode = "00"
	c.res.Data = "08D7B4FB629D0885H1.25.11M1.17.02C1.16.10V1310-000035"
	return nil
}
