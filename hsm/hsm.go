package hsm

import (
	"bufio"
	"encoding/binary"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
)

type Hsm struct {
	conn net.Conn
}

func NewHsm(conn net.Conn) *Hsm {
	return &Hsm{conn}
}

func (hsm *Hsm) Read() ([]byte, error) {
	reader := bufio.NewReader(hsm.conn)

	var msgSize int16
	// message size
	err := binary.Read(reader, binary.BigEndian, &msgSize)
	if err != nil {
		return nil, err
	}
	log.Printf("读取数据头字节长度=%d", msgSize)
	//从缓存区读取大小为数据长度的数据
	data := make([]byte, msgSize)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}
	log.Printf("读取数据=%X", data)
	return data, nil
}

func (hsm *Hsm) Write(in []byte) {
	outBytes := LengthFieldPrepend(in)
	hsm.conn.Write(outBytes)
}

// 填充数据长度
func LengthFieldPrepend(in []byte) []byte {
	dataLength := len(in)
	lenSize := 2
	out := make([]byte, dataLength+lenSize)
	lenBytes := make([]byte, lenSize)
	binary.BigEndian.PutUint16(lenBytes, uint16(dataLength))
	out[0] = lenBytes[0]
	out[1] = lenBytes[1]
	for i := 0; i < dataLength; i++ {
		out[lenSize+i] = in[i]
	}
	return out
}
