package biz

import (
	"bytes"
	"context"
	"crypto/md5"
	"ehc-hsm-server/utils"
	"encoding/binary"
	"encoding/hex"
	"errors"
	log "github.com/sirupsen/logrus"
)

type BizHandle interface {
	//业务处理
	Handle() error

	//解码
	Decode(buf []byte) error

	//编码
	Encode() ([]byte, error)
}

const (
	REQ_CODE_NC        = "NC" //获取设备信息
	REQ_CODE_CR        = "CR" //心跳
	REQ_CODE_DIGEST_3C = "3C" //数据加解密
	REQ_CODE_S3        = "S3" //加密（生成健康卡、加密时间）
	REQ_CODE_S4        = "S4" //解密（验证健康卡、解密时间）
)

// 安全密钥
const SAFE_KEY string = "0123456789ABCDEFFEDCBA9876543210"

// 保护密钥
const PROTECT_KEY string = "11223344556677888877665544332211"

type HsmReq struct {
	ReqCode string
}

/**
 * 解码
 *
 * @param in 字节缓存
 */
func (*HsmReq) decode(in bytes.Buffer) {

}

var ERROR_REQ = errors.New("Error reqCode")

func HandMsg(ctx context.Context, buf []byte) []byte {
	log.WithField("tradeId", ctx.Value("tradeId")).Printf("请求数据=[%s]", hex.EncodeToString(buf))

	in := bytes.NewBuffer(buf)
	reqCode, err := utils.ReadString(in, 2)
	if err != nil {
		return nil
	}

	var out []byte
	var biz BizHandle
	if REQ_CODE_NC == reqCode {
		//获取设备信息
		biz = NewNc()
	} else if REQ_CODE_CR == reqCode {
		//心跳
		biz = NewCr()
	} else if REQ_CODE_DIGEST_3C == reqCode {
		//3c生成摘要
		biz = NewDigest3c()
	} else if REQ_CODE_S3 == reqCode {
		//生成健康卡或加密时间
		biz = NewS3()
	} else if REQ_CODE_S4 == reqCode {
		//验证健康卡或解密时间
		biz = NewS4()
	}

	biz.Decode(buf)
	biz.Handle()
	out, err = biz.Encode()

	hexStr := hex.EncodeToString(out)
	log.WithField("tradeId", ctx.Value("tradeId")).Printf("响应数据=[%X]", out)
	dataLength := len(hexStr)
	log.WithField("tradeId", ctx.Value("tradeId")).Printf("响应数据=%d,%s", dataLength, hexStr)

	log.Printf("响应数据md5=%X", md5.Sum(out))
	log.Printf("响应数据=[%s]", hex.EncodeToString(out))
	out = LengthFieldPrepend(out)
	return out
}

// 填充数据长度
func LengthFieldPrepend(in []byte) []byte {
	dataLength := len(in)
	lenSize := 2
	out := make([]byte, dataLength+lenSize)
	lenBytes := make([]byte, lenSize)
	binary.BigEndian.PutUint16(lenBytes, uint16(dataLength))
	out[0] = lenBytes[0]
	out[1] = lenBytes[1]
	for i := 0; i < dataLength; i++ {
		out[lenSize+i] = in[i]
	}
	return out
}
