package main

import (
	"bufio"
	"context"
	"ehc-hsm-server/biz"
	"encoding/binary"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"strings"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	address := "127.0.0.1:18018"
	log.Println("run hsm server in", address)
	listen, err := net.Listen("tcp", address)
	checkError(err)
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		log.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		go handConn(conn)
	}
}

func handConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	//writer := bufio.NewWriter(conn)
	for {
		var msgSize int16
		// message size
		err := binary.Read(reader, binary.BigEndian, &msgSize)
		if err != nil {
			return
		}

		uuid := strings.ReplaceAll(uuid.New().String(), "-", "")
		ctx := context.WithValue(context.Background(), "tradeId", uuid)
		log.Printf("读取数据头字节长度=%d,\n", msgSize)
		//从缓存区读取大小为数据长度的数据
		data := make([]byte, msgSize)
		_, err = io.ReadFull(reader, data)
		if err != nil {
			continue
		}

		outBytes := biz.HandMsg(ctx, data)
		conn.Write(outBytes)
	}
}

func checkError(err error) {
	if err != nil {
		log.Error("Fatal error ", err.Error())
		os.Exit(1)
	}
}
