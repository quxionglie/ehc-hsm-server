package utils

import (
	"bytes"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"testing"
)

var src = "004a533330303030414b303130343031424534423341313334383730313538393431423443354543423738464541373630303032303031343031343430333033313939343035313039303131"

func TestReadBytes(t *testing.T) {
	byte, err := hex.DecodeString(src)
	if err != nil {
		t.Fail()
	}

	buf := bytes.NewBuffer(byte)
	out, err := ReadBytes(buf, 2)
	if err != nil {
		t.Fail()
	}
	log.Println("ReadBytes=", hex.EncodeToString(out))
}

func TestReadInt(t *testing.T) {
	byte, err := hex.DecodeString("0035")
	if err != nil {
		t.Fail()
	}

	buf := bytes.NewBuffer(byte)
	out, err := ReadInt(buf, 2)
	if err != nil {
		t.Fail()
	}
	log.Println("ReadInt=", out)
}

func TestReadIntHex(t *testing.T) {
	byte, err := hex.DecodeString("0035")
	if err != nil {
		t.Fail()
	}

	buf := bytes.NewBuffer(byte)
	out, err := ReadIntHex(buf, 2)
	if err != nil {
		t.Fail()
	}
	log.Println("ReadIntHex2=", out)
}

//
//func TestReadString(t *testing.T) {
//	res := biz.NewNcRes()
//	res.ErrCode = "00"
//	res.Data = "08D7B4FB629D0885H1.25.11M1.17.02C1.16.10V1310-000035"
//	out := res.Encode()
//	log.Println(hex.EncodeToString(out)) // 104
//}
