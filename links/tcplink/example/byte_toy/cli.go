package main

import (
	"fmt"
	"net"
	"sync"

	"time"

	"chargerlink.com/golib/link/tcplink"
	"chargerlink.com/golib/link/tcplink/codec"
)

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	for i := 0; i < 30000; i++ {
		go clientSession()
	}
	wg.Wait()
}

func clientSession() {
	bt := codec.Byte(3, 1, 3, 0, 0)
	session, _, err := tcplink.ConnectTimeout("tcp", "127.0.0.1:59999", 30*time.Second, bt, 0)
	if err != nil {
		fmt.Printf("err:[%s]\r\n", err.Error())
		return
	}

	buff := []byte{0x68, 0x08, 0x00, 0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	go func() {
		for i := 0; i < 1000; i++ {
			if err = session.Send(buff); err != nil {
				fmt.Printf("send err:[%s]\r\n", err.Error())
				return
			}
			time.Sleep(5 * time.Second)
		}
	}()

	for {
		data, err := session.Receive()
		if err != nil {
			fmt.Printf("recv err:[%s]\r\n", err.Error())
			return
		}
		fmt.Printf("recv and recv data:[% x] \r\n", data)
	}
}

func client() {
	conn, err := net.Dial("tcp", "127.0.0.1:59999")
	if err != nil {
		fmt.Printf("err:[%s]\r\n", err.Error())
		return
	}
	defer conn.Close()

	buff := []byte{0x68, 0x08, 0x00, 0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	_, err = conn.Write(buff)
	if err != nil {
		fmt.Printf("write err:[%s]\r\n", err.Error())
		return
	}

	rbuf := make([]byte, 1024)
	_, err = conn.Read(rbuf)
	if err != nil {
		fmt.Printf("read err:[%s]\r\n", err.Error())
		return
	}
	fmt.Printf("send and recv data:[% x] \r\n", rbuf)
}
