package utils

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEbcSm4(t *testing.T) {
	srcStr := "92CDCE45A16799FD6D3231BA5E986602"
	keyHex := "0123456789ABCDEFFEDCBA9876543210"
	key, _ := hex.DecodeString(keyHex)
	enc, _ := Sm4EncryptNoPadding(key, []byte(srcStr))

	stdOut := "b3cd4babf4a40ec2065bf4783a3a9a7fc57599135f9f7f0fca60d8639b0a025f"

	dec, _ := Sm4DecryptNoPadding(key, enc)
	fmt.Printf("原始原串: %s\n", srcStr)
	fmt.Printf("解密原串: %s\n", dec)
	fmt.Printf("密钥: %s\n", keyHex)
	fmt.Printf("计算加密: %x\n", enc)
	fmt.Printf("标准加密: %s\n", stdOut)
	assert.Equal(t, hex.EncodeToString(enc), stdOut)

}
