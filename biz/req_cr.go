package biz

import (
	"bytes"
	"ehc-hsm-server/pkg/rwbytes"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strings"
)

//心跳
type CrReq struct {
	ReqCode string
	Seq     string
}

type CrRes struct {
	ReqCode string
	ErrCode string
	Data    string
}
type Cr struct {
	req *CrReq
	res *CrRes
}

func NewCr() *Cr {
	req := CrReq{}
	res := CrRes{ReqCode: "CS", ErrCode: "00"}
	return &Cr{&req, &res}
}
func (req *CrReq) Decode(in []byte) error {
	buf := bytes.NewBuffer(in)
	req.ReqCode, _ = rwbytes.ReadString(buf, 2)
	req.Seq, _ = rwbytes.ReadString(buf, 4)
	log.Printf("请求参数：%v", req)
	return nil
}

func (req *CrReq) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	rwbytes.WriteString(buf, 2, req.ReqCode)
	rwbytes.WriteString(buf, 4, req.Seq)
	return buf.Bytes(), nil
}

func (res *CrRes) Decode(in []byte) error {
	buf := bytes.NewBuffer(in)
	res.ReqCode, _ = rwbytes.ReadString(buf, 2)
	res.ErrCode, _ = rwbytes.ReadString(buf, 2)
	if len(in) >= 20 {
		//Res.Data, _ = rwbytes.ReadString(buf, 16)
	}
	return nil
}

func (res *CrRes) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write([]byte(res.ReqCode))
	buf.Write([]byte(res.ErrCode))
	if res.Data != "" {
		buf.Write([]byte(res.Data))
	}
	return buf.Bytes(), nil
}

func (c *Cr) Decode(buf []byte) error {
	return c.req.Decode(buf)
}
func (c *Cr) Encode() ([]byte, error) {
	return c.res.Encode()
}

func (c *Cr) Handle() error {
	c.res.ErrCode = "00"
	//java: crRes.setData(UUID.randomUUID().toString().replaceAll("-", "").substring(0, 16));
	uuid := strings.ReplaceAll(uuid.New().String(), "-", "")
	c.res.Data = uuid[0:16]
	return nil

}
