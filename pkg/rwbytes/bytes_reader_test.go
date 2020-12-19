package rwbytes

import (
	"bytes"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "004a", hex.EncodeToString(out))
}

func TestReadInt(t *testing.T) {
	buf := bytes.NewBufferString("35")
	out, err := ReadInt(buf, 2)
	if err != nil {
		t.Fail()
	}
	log.Printf("ReadInt=%d", out)
	assert.Equal(t, int32(35), out)
}

func TestReadIntHex(t *testing.T) {
	buf := bytes.NewBufferString("35")
	out, err := ReadIntHex(buf, 2)
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, int32(53), out)

	out, err = ReadIntHex(bytes.NewBufferString("0A"), 2)
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, int32(10), out)
}
