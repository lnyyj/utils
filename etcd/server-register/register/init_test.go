package srvregister

import (
	"context"
	"testing"
)

func Test_reg(t *testing.T) {
	{
		sr := New(`{"endpoints":["localhost:2379"]}`)
		sr.Reg(context.TODO(), "/user/node1", "10.10.2.1:57777")
	}
	{
		sr := New(`{"endpoints":["localhost:2379"]}`)
		sr.Reg(context.TODO(), "/msg/node1", "10.10.2.1:56666")
	}

	select {}
}
