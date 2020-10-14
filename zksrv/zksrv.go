package zksrv

import (
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

// ZKSrv zk服务公共对象
type zkSrv struct {
	zkConn *zk.Conn
	hosts  []string
	stop   chan bool // 关闭标识
}

func (zks *zkSrv) checkzk() {
	if zks.zkConn == nil {
		panic("please create zookeeper client")
	}
}

// CreateZKCli 创建zk客户端
func (zks *zkSrv) CreateZKCli() (err error) {
	zks.zkConn, _, err = zk.Connect(zks.GetHosts(), time.Second*5)
	if err != nil {
		return err
	}
	zks.stop = make(chan bool)
	return nil
}

// StopZKCli 关闭TCP连接
func (zks *zkSrv) StopZKCli() {
	zks.zkConn.Close()
	zks.stop <- true
}

// GetZKCli 获取zk客户端
func (zks *zkSrv) GetZKCli() *zk.Conn {
	return zks.zkConn
}

// SetHosts 是设置zookeeper服务器连接地址
func (zks *zkSrv) SetHosts(host ...string) *zkSrv {
	for _, v := range host {
		zks.hosts = append(zks.hosts, v)
	}
	return zks
}

// GetHosts 是获取zookeeper服务器连接地址
func (zks *zkSrv) GetHosts(host ...string) []string {
	return zks.hosts
}
