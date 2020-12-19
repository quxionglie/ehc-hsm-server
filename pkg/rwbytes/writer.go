package rwbytes

import (
	"bytes"
	"fmt"
	"strconv"
)

func toIso8859_1(data string) []byte {
	return []byte(data)
}

/**
 * 写字符串
 */
func WriteString(in *bytes.Buffer, fixLen int, data string) (int, error) {
	out := toIso8859_1(data)
	return WriteBytes(in, fixLen, out)
}

/**
 * 写数组
 *
 * @param in ByteBuf
 * @param fixLen 长度
 * @return 字符数组
 */
func WriteBytes(in *bytes.Buffer, fixLen int, datas []byte) (int, error) {
	putBytes := make([]byte, fixLen)
	for i := range datas {
		putBytes[i] = datas[i]
	}
	return in.Write(putBytes)
}

/**
 * 写整数
 */
func WriteInt(in *bytes.Buffer, fixLen int, data int) (int, error) {
	lenStr := fmt.Sprintf("%0"+strconv.Itoa(fixLen)+"d", data)
	return in.Write([]byte(lenStr))
}

/**
 * 读整数
 * @param in ByteBuf
 * @param length 长度
 * @return 整数
 */
func WriteIntHex(in *bytes.Buffer, fixLen int, data int) (int, error) {
	lenStr := fmt.Sprintf("%0"+strconv.Itoa(fixLen)+"X", data)
	return in.Write([]byte(lenStr))
}
