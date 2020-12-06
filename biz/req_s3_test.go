package biz

import (
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestS3decode(t *testing.T) {
	//30303030414B30313034303135464533314433454234434344314130314345324331344233333235324530303032303031343031313130313031313939323131313430363139
	in := "004a533330303030414b303130343031374244363737414243463646353738343834323938383534333039304138374230303032303031343031313130313031313939323131313430363139"
	byte, err := hex.DecodeString(in)
	if err != nil {
		t.Fail()
	}

	//buf := bytes.NewBuffer(byte[2:])
	s3 := NewS3()
	s3.Decode(byte[2:])
	s3.Handle()
	out, _ := s3.Encode()
	stdOut := strings.ToUpper("0028533430303030323037b098fc786ee122f2a8d4a1cb4b8dee382e7813958914f15d3426f89e2e20dd")
	outHexStr := strings.ToUpper(hex.EncodeToString(LengthFieldPrepend(out)))
	log.Printf("标准数据=[%s]", stdOut)
	log.Printf("响应数据=[%s]", outHexStr)
	assert.Equal(t, stdOut, outHexStr)
}
