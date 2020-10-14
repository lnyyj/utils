package thriftlink

import (
	"fmt"
	"os"

	"time"

	"reflect"

	"git.apache.org/thrift.git/lib/go/thrift"
)

// 接收收据回调函数对象
type Instance interface {
	NewInstance(thrift.TTransport, thrift.TProtocolFactory) interface{}
}

// thrift 客户端封装
type ThriftCli struct {
	addr string

	transport *thrift.TSocket
}

func NewThriftCli() *ThriftCli {
	return &ThriftCli{}
}

func (t *ThriftCli) SetServerAddr(addr string) *ThriftCli {
	t.addr = addr
	return t
}

func (t *ThriftCli) EnableConn(fun interface{}) (instance interface{}, err error) {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	t.transport, err = thrift.NewTSocket(t.addr)
	t.transport.SetTimeout(15 * time.Second)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error resolving address:", err)
		return nil, err
	}
	useTransport := transportFactory.GetTransport(t.transport)

	params := make([]reflect.Value, 2)
	params[0] = reflect.ValueOf(useTransport).Elem()
	params[1] = reflect.ValueOf(protocolFactory).Elem()

	instance = reflect.ValueOf(fun).Call(params)
	if err = t.transport.Open(); err != nil {
		return nil, err
	}
	return instance, nil
}

func (t *ThriftCli) IsConn() bool {
	return t.transport.IsOpen()
}

func (t *ThriftCli) DeableConn() {
	t.transport.Close()
}
