package biz

import (
	"bytes"
	"ehc-hsm-server/pkg/rwbytes"
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
	Req *NcReq
	Res *NcRes
}

func NewNc() *Nc {
	req := NcReq{REQ_CODE_NC}
	res := NcRes{ReqCode: "ND", ErrCode: "00"}
	return &Nc{&req, &res}
}

func (req *NcReq) Decode(in []byte) (err error) {
	buf := bytes.NewBuffer(in)
	req.ReqCode, err = rwbytes.ReadString(buf, 2)
	return
}

func (req *NcReq) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write([]byte(REQ_CODE_NC))
	return buf.Bytes(), nil
}

func (res *NcRes) Decode(in []byte) error {
	buf := bytes.NewBuffer(in)
	res.ReqCode, _ = rwbytes.ReadString(buf, 2)
	res.ErrCode, _ = rwbytes.ReadString(buf, 2)
	if len(in) >= 20 {
		//Res.Data, _ = rwbytes.ReadString(buf, 16)
	}
	return nil
}

func (res *NcRes) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write([]byte(res.ReqCode))
	buf.Write([]byte(res.ErrCode))
	if res.Data != "" {
		buf.Write([]byte(res.Data))
	}
	return buf.Bytes(), nil
}

func (c *Nc) Decode(buf []byte) error {
	return c.Req.Decode(buf)
}
func (c *Nc) Encode() ([]byte, error) {
	return c.Res.Encode()
}

func (c *Nc) Handle() error {
	//Req := NewNcReq()
	//Req.Decode(buf)
	//log.Printf("请求参数：%v", Req)
	c.Res.ErrCode = "00"
	c.Res.Data = "08D7B4FB629D0885H1.25.11M1.17.02C1.16.10V1310-000035"
	return nil
}
