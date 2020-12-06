package utils

import (
	"bytes"
	"strconv"
)

func toUtf8(iso8859_1_buf []byte) string {
	buf := make([]rune, len(iso8859_1_buf))
	for i, b := range iso8859_1_buf {
		buf[i] = rune(b)
	}
	return string(buf)
}

/**
 * 读字符串
 *
 * @param in     ByteBuf
 * @param length 长度
 * @return 字符串
 */
func ReadString(in *bytes.Buffer, length int32) (string, error) {
	bytes, err := ReadBytes(in, length)
	if err != nil {
		return "", err
	}
	//return new String(bytes, Charsets.ISO_8859_1)
	return toUtf8(bytes), nil
}

/**
 * 读字符数组
 *
 * @param in     ByteBuf
 * @param length 长度
 * @return 字符数组
 */
func ReadBytes(in *bytes.Buffer, length int32) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := in.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

/**
 * 读整数
 *
 * @param in     ByteBuf
 * @param length 长度
 * @return 整数
 */
func ReadInt(in *bytes.Buffer, length int32) (int32, error) {
	s, err := ReadString(in, length)
	if err != nil {
		return 0, err
	}
	// string 转 int32
	j, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(j), nil
}

/**
 * 读整数
 * @param in ByteBuf
 * @param length 长度
 * @return 整数
 */
func ReadIntHex(in *bytes.Buffer, length int32) (int32, error) {
	s, err := ReadString(in, length)
	if err != nil {
		return 0, err
	}
	j, err := strconv.ParseInt(s, 16, 32)
	if err != nil {
		return 0, err
	}
	return int32(j), nil
}
