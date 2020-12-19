package biz

import (
	"bytes"
	"ehc-hsm-server/pkg/rwbytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tjfoc/gmsm/sm3"
)

//3c生成摘要
type Digest3cReq struct {
	ReqCode          string
	AlgorithmPattern string
	DataLength       int
	Data             []byte
	EndStr           string
	Var3Length       int
	Var4Length       int
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

func (req *Digest3cReq) Decode(buf []byte) error {
	in := bytes.NewBuffer(buf)

	var err error = nil
	req.ReqCode, err = rwbytes.ReadString(in, 2)
	req.AlgorithmPattern, err = rwbytes.ReadString(in, 2)
	req.DataLength, err = rwbytes.ReadInt(in, 4)
	req.Data, err = rwbytes.ReadBytes(in, req.DataLength)
	req.EndStr, err = rwbytes.ReadString(in, 1)
	req.Var3Length, err = rwbytes.ReadInt(in, 4)
	req.Var4Length, err = rwbytes.ReadInt(in, 4)
	log.Printf("请求参数：%v", req)
	return err
}

func (req *Digest3cReq) Encode() ([]byte, error) {
	in := new(bytes.Buffer)
	rwbytes.WriteString(in, 2, req.ReqCode)
	rwbytes.WriteString(in, 2, req.AlgorithmPattern)
	rwbytes.WriteInt(in, 4, req.DataLength)
	rwbytes.WriteBytes(in, req.DataLength, req.Data)
	rwbytes.WriteString(in, 1, req.EndStr)
	rwbytes.WriteInt(in, 4, req.Var3Length)
	rwbytes.WriteInt(in, 4, req.Var4Length)
	return in.Bytes(), nil
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

func (res *Digest3cRes) Encode() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write([]byte(res.ReqCode))
	buf.Write([]byte(res.ErrCode))

	//byte[] len = Strings.padStart(DataLength + "", 2, '0').getBytes(Charsets.ISO_8859_1);
	lenStr := fmt.Sprintf("%02d", res.DataLength)
	buf.Write([]byte(lenStr))
	if res.Data != nil {
		buf.Write(res.Data)
	}
	return buf.Bytes(), nil
}

func (res *Digest3cRes) Decode() ([]byte, error) {
	buf := bytes.Buffer{}
	buf.Write([]byte(res.ReqCode))
	buf.Write([]byte(res.ErrCode))

	lenStr := fmt.Sprintf("%02d", res.DataLength)
	buf.Write([]byte(lenStr))
	if res.Data != nil {
		buf.Write(res.Data)
	}
	return buf.Bytes(), nil
}

func (c *Digest3c) Decode(buf []byte) error {
	return c.req.Decode(buf)
}
func (c *Digest3c) Encode() ([]byte, error) {
	return c.res.Encode()
}

func (c *Digest3c) Handle() error {
	c.res.ErrCode = "00"
	h := sm3.New()
	h.Write(c.req.Data)
	sum := h.Sum(nil)
	c.setData(sum)
	return nil
}
