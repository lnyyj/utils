package tcplink

import (
	"net"
	"time"
)

type Context interface{}

// 协议接口
type Protocol interface {
	NewCodec(conn net.Conn) (Codec, Context, error)
}

// 协议回调函数
type ProtocolFunc func(conn net.Conn) (Codec, Context, error)

func (pf ProtocolFunc) NewCodec(conn net.Conn) (Codec, Context, error) {
	return pf(conn)
}

// 编码接口
type Codec interface {
	Receive() (interface{}, error)
	Send(interface{}) error
	Close() error
}

// 监听地址并创建一个服务
func Serve(network, address string, protocol Protocol, sendChanSize int, mon bool) (*Server, error) {
	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP(network, tcpAddr)
	//listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}
	return NewServer(listener, protocol, sendChanSize, mon), nil
}

// 客户端连接并创建一个会话
func Connect(network, address string, protocol Protocol, sendChanSize int) (*Session, Context, error) {
	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return nil, nil, err
	}

	// conn, err := net.Dial(network, address)
	conn, err := net.DialTCP(network, nil, tcpAddr)
	if err != nil {
		return nil, nil, err
	}
	codec, ctx, err := protocol.NewCodec(conn)
	if err != nil {
		return nil, nil, err
	}
	return NewSession(codec, sendChanSize), ctx, nil
}

// 客户端超时连接并创建一个会话
func ConnectTimeout(network, address string, timeout time.Duration, protocol Protocol, sendChanSize int) (*Session, Context, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, nil, err
	}
	codec, ctx, err := protocol.NewCodec(conn)
	if err != nil {
		return nil, nil, err
	}
	return NewSession(codec, sendChanSize), ctx, nil
}
