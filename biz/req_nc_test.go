package biz

import (
	"crypto/md5"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestNc2(t *testing.T) {
	nc := NewNc()
	nc.Decode([]byte("NC"))
	nc.Handle()
	res, _ := nc.Encode()
	toString := hex.EncodeToString(res)
	log.Printf("%d,%s", len(toString), toString) // 104

	//310A76F52C14F5FFE333D5507116294E
	log.Printf("md5=%X", md5.Sum(res))

	n := len(res)
	bytes := res[0:n]
	log.Printf("md5=%X", md5.Sum(bytes))
	//310A76F52C14F5FFE333D5507116294E
}
