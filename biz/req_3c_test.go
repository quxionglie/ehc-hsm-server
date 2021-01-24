package biz

import (
	"ehc-hsm-server/hsm"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test3c(t *testing.T) {
	//3C:002e33433230303032393031313130313031313939323131313430363139e8b4b5e69c8de9809a3b3030303030303030
	//3D:0026334430303332a59e97c486e8f1c29f46ea87a515e61db1eecf23076b4993f0e0c5cbebf909c8
	src := "002e33433230303032393031313130313031313939323131313430363139e8b4b5e69c8de9809a3b3030303030303030"
	byte, err := hex.DecodeString(src)
	if err != nil {
		t.Fail()
	}

	req := NewDigest3c()
	req.Decode(byte[2:])
	req.Handle()
	req.Encode()

	//buf := bytes.NewBuffer(byte[2:])
	out, _ := req.Encode()
	stdOut := strings.ToUpper("0026334430303332a59e97c486e8f1c29f46ea87a515e61db1eecf23076b4993f0e0c5cbebf909c8")
	outHexStr := strings.ToUpper(hex.EncodeToString(hsm.LengthFieldPrepend(out)))
	log.Printf("标准数据=[%s]", stdOut)
	log.Printf("响应数据=[%s]", outHexStr)
	assert.Equal(t, stdOut, outHexStr)
}
