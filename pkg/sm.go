package pkg

import (
	"bytes"
	"encoding/binary"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
)

func Sm3(data string) []byte {
	h := sm3.New()
	h.Write([]byte(data))
	sum := h.Sum(nil)
	return sum
}

func Sm4EncryptNoPadding(key []byte, plainText []byte) ([]byte, error) {
	return sm4EcbNoPadding(key, plainText, true)
}

func Sm4DecryptNoPadding(key []byte, cipherText []byte) ([]byte, error) {
	return sm4EcbNoPadding(key, cipherText, false)
}

//mode:true加密，faLse解密
func sm4EcbNoPadding(key []byte, in []byte, mode bool) (out []byte, err error) {
	var inData []byte = in
	out = make([]byte, len(inData))
	c, err := sm4.NewCipher(key)
	if err != nil {
		panic(err)
	}
	if mode {
		for i := 0; i < len(inData)/16; i++ {
			in_tmp := inData[i*16 : i*16+16]
			out_tmp := make([]byte, 16)
			c.Encrypt(out_tmp, in_tmp)
			copy(out[i*16:i*16+16], out_tmp)
		}
	} else {
		for i := 0; i < len(inData)/16; i++ {
			in_tmp := inData[i*16 : i*16+16]
			out_tmp := make([]byte, 16)
			c.Decrypt(out_tmp, in_tmp)
			copy(out[i*16:i*16+16], out_tmp)
		}
	}

	return out, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding) //用0去填充
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

// pkcs5填充
func pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	if length == 0 {
		return nil
	}
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
