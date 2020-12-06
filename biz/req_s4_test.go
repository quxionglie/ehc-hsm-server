package biz

import (
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestS4(t *testing.T) {
	in := "0056533430303030414b3031303430313742443637374142434636463537383438343239383835343330393041383742303030353030323037b098fc786ee122f2a8d4a1cb4b8dee382e7813958914f15d3426f89e2e20dd"
	byte, err := hex.DecodeString(in)
	if err != nil {
		t.Fail()
	}

	biz := NewS4()
	biz.Decode(byte[2:])
	biz.Handle()
	//buf := bytes.NewBuffer(byte[2:])
	out, err := biz.Encode()
	stdOut := strings.ToUpper("002853353030303032303031313130313031313939323131313430363139000000000000000000000000")

	outHexStr := strings.ToUpper(hex.EncodeToString(LengthFieldPrepend(out)))
	log.Printf("标准数据=[%s]", stdOut)
	log.Printf("响应数据=[%s]", outHexStr)
	assert.Equal(t, stdOut, outHexStr)
}
