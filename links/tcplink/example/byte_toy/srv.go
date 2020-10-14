package main

import (
	"fmt"

	"chargerlink.com/golib/idgen"
	"chargerlink.com/golib/link/tcplink"
	"chargerlink.com/golib/link/tcplink/codec"
	"chargerlink.com/golib/log"
)

func main() {
	fmt.Println("------> test byte byte toy begin <------")
	log.InitConsoleOutput(7)
	// server()
	server2()
}

func server() {
	bt := codec.Byte(3, 1, 3, 60, 0)
	server, err := tcplink.Serve("tcp", ":59999", bt, 0, true)
	if err != nil {
		fmt.Println("create service fial:", err)
		return
	}

	server.Serve(tcplink.HandlerFunc(sessionLoop))
}

func server2() {
	bt := codec.Byte(3, 1, 3, 60, 0)
	server, err := tcplink.Serve("tcp", ":59999", bt, 0, true)
	if err != nil {
		fmt.Println("create service fial:", err)
		return
	}

	server.ServeEx(tcplink.HandlerFunc(sessionLoop), tcplink.InitSessionFunc(SessionInit))
}

func SessionInit(codec tcplink.Codec) (interface{}, uint64, error) {
	data, err := codec.Receive()
	if err != nil {
		fmt.Println("init session err")
		return nil, 0, err
	}
	id, err := idgen.CreateId()

	return data.([]byte), uint64(id), err
}

func sessionLoop(session *tcplink.Session, _ interface{}, sessionErr error) {
	if sessionErr != nil {
		fmt.Printf("session loop error:[%s]\r\n", sessionErr.Error())
		return
	}

	buff := []byte{0x68, 0x08, 0x00, 0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	session.Send(buff)
	session.Send(buff)
	session.Send(buff)
	session.Send(buff)
	for {
		// fmt.Println("---------> start receive")
		req, err := session.Receive()
		if err != nil {
			fmt.Println("receive err: ", err)
			return
		}

		// fmt.Printf("--->1 recv and send data:[% x] \r\n", req.([]byte))

		err = session.Send(req)
		if err != nil {
			fmt.Println("send error: ", err)
			return
		}
		// fmt.Printf("---> recv and send data:[% x] \r\n", req.([]byte))
	}
}
