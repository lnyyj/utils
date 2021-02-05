package srvregister

import (
	"context"
	"encoding/json"
	"log"

	"go.etcd.io/etcd/clientv3"
)

// ServerRegister 服务注册
// todo: 支持注册多个服务
type ServerRegister struct {
	Config clientv3.Config
	cli    *clientv3.Client

	lease *clientv3.LeaseGrantResponse
}

// New ...
func New(config string) *ServerRegister {
	ret := &ServerRegister{}
	if config != "" {
		if err := json.Unmarshal([]byte(config), &ret.Config); err != nil {
			panic(err.Error())
		}
	}
	return ret
}

// Reg  服务注册
func (sr *ServerRegister) Reg(ctx context.Context, name, addr string) (err error) {
	if sr.cli, err = clientv3.New(sr.Config); err != nil {
		return
	} else if sr.lease, err = sr.cli.Lease.Grant(ctx, 10); err != nil {
		return
	} else if _, err = sr.cli.KV.Put(ctx, name, addr, clientv3.WithLease(sr.lease.ID)); err != nil {
		return
	}

	if leaseRespChan, err := sr.cli.KeepAlive(ctx, sr.lease.ID); err != nil {
		return err
	} else {
		go func(lka <-chan *clientv3.LeaseKeepAliveResponse) {
			for leaseKeepResp := range lka {
				log.Println("续约成功", leaseKeepResp)
			}
		}(leaseRespChan)
	}
	return
}

// Close 关闭
func (sr *ServerRegister) Close(ctx context.Context) (err error) {
	if _, err = sr.cli.Lease.Revoke(ctx, sr.lease.ID); err != nil {
		return
	}
	return sr.cli.Close()
}
