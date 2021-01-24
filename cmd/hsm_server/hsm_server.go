package main

import (
	"ehc-hsm-server/biz"
	"ehc-hsm-server/hsm"
	log "github.com/sirupsen/logrus"
	"net"
	"os"
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
	address := ":18018"
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
	for {
		hsm := hsm.NewHsm(conn)

		data, err := hsm.Read()
		if err != nil {
			continue
		}

		outBytes := biz.HandMsg(data)
		hsm.Write(outBytes)
	}
}

func checkError(err error) {
	if err != nil {
		log.Error("Fatal error ", err.Error())
		os.Exit(1)
	}
}
