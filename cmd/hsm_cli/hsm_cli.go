package main

import (
	"context"
	"ehc-hsm-server/biz"
	"ehc-hsm-server/hsm"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net"
	"os"
)

func main() {
	address := "127.0.0.1:18018"
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		return test3c(conn)
	})

	err = g.Wait()
	if err != nil {
		fmt.Println("Error run:", err)
		os.Exit(1)
	}
}

func testNc(conn net.Conn) error {
	h := hsm.NewHsm(conn)
	nc := biz.NewNc()
	out, _ := nc.Req.Encode()
	h.Write(out)
	h.Read()
	return nil
}

//3c生成摘要
func test3c(conn net.Conn) error {
	h := hsm.NewHsm(conn)
	req := &biz.Digest3cReq{}
	req.ReqCode = biz.REQ_CODE_DIGEST_3C
	req.AlgorithmPattern = "20"
	req.Data = []byte("01110101199211140619李三三")
	req.DataLength = len(req.Data)
	req.EndStr = ";"
	req.Var3Length = 0
	req.Var4Length = 0
	out, _ := req.Encode()
	h.Write(out)
	h.Read()
	return nil
}
