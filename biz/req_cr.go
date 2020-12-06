package biz

import (
	"bytes"
	"ehc-hsm-server/utils"
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

func (c *Cr) Decode(in []byte) error {
	buf := bytes.NewBuffer(in)
	c.req.ReqCode, _ = utils.ReadString(buf, 2)
	c.req.Seq, _ = utils.ReadString(buf, 4)
	log.Printf("请求参数：%v", c.req)
	return nil
}

func (c *Cr) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write([]byte(c.res.ReqCode))
	buf.Write([]byte(c.res.ErrCode))
	if c.res.Data != "" {
		buf.Write([]byte(c.res.Data))
	}
	return buf.Bytes(), nil
}

func (c *Cr) Handle() error {
	c.res.ErrCode = "00"
	//java: crRes.setData(UUID.randomUUID().toString().replaceAll("-", "").substring(0, 16));
	uuid := strings.ReplaceAll(uuid.New().String(), "-", "")
	c.res.Data = uuid[0:16]
	return nil

}
