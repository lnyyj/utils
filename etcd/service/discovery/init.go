package srvdiscovery

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"go.etcd.io/etcd/clientv3"
)

var (
	watchServers map[string]*server

	// Config 连接etcd配置
	Config clientv3.Config
	cli    *clientv3.Client

	lock sync.Locker
)

type server struct {
	cur, count   int // 获取时使用
	name, prefix string
	nodes        []*node
}

type node struct {
	key, addr string // 地址信息
	load      int    // 计算下负载信息 todo: 下个版本执行
}

func init() {
	watchServers = make(map[string]*server)
}

// Add ...
func Add(srvName, prefix string) {
	sinfo, ok := watchServers[srvName]
	if !ok {
		sinfo = &server{}
	}
	sinfo.name, sinfo.prefix = srvName, prefix
	watchServers[srvName] = sinfo
}

// Watch watch server keys
func Watch(config string) (err error) {
	if config != "" {
		if err = json.Unmarshal([]byte(config), Config); err != nil {
			return
		}
	}
	if cli, err = clientv3.New(Config); err != nil {
		return
	}

	// 初始化
	for _, server := range watchServers {
		resp, err := cli.Get(context.TODO(), server.prefix, clientv3.WithPrefix())
		if err != nil {
			return err
		}
		for _, ev := range resp.Kvs {
			addServiceList(server.name, string(ev.Key), string(ev.Value))
		}
	}

	for name, server := range watchServers {
		go watch(name, server.prefix)
	}
	return
}

func watch(name, prefix string) {
	wch := cli.Watch(context.TODO(), prefix, clientv3.WithPrefix())

	for watchResponse := range wch {
		// if watchResponse.Canceled || watchResponse.Err() != nil { // 取消watch
		// 	break
		// }
		for _, event := range watchResponse.Events {
			if event.IsCreate() { // 新增
				addServiceList(name, string(event.Kv.Key), string(event.Kv.Value))
			} else if event.IsModify() {
				modifyServiceList(name, string(event.Kv.Key), string(event.Kv.Value))
			} else if event.Type == clientv3.EventTypeDelete {
				delServiceList(name, string(event.Kv.Key))
			}
		}
	}

}

//addServiceList 新增服务地址
func addServiceList(name, key, val string) {
	lock.Lock()
	defer lock.Unlock()

	srvInfo, _ := watchServers[name]
	srvInfo.nodes = append(srvInfo.nodes, &node{key: key, addr: val})
	srvInfo.count++
	log.Println("put key :", key, "val:", val)
}

// modifyServiceList 新增服务地址
func modifyServiceList(name, key, val string) {
	lock.Lock()
	defer lock.Unlock()

	srvInfo, _ := watchServers[name]
	for _, v := range srvInfo.nodes {
		if v.key == key {
			v.addr = val
			break
		}
	}
	log.Println("change key :", key, "val:", val)
}

//delServiceList 删除服务地址
func delServiceList(name, key string) {
	lock.Lock()
	defer lock.Unlock()
	srvInfo, _ := watchServers[name]
	for k := range srvInfo.nodes {
		if srvInfo.nodes[k].key == key {
			srvInfo.count--
			srvInfo.nodes = append(srvInfo.nodes[:k], srvInfo.nodes[k+1:]...)
			break
		}
	}
	log.Println("del key:", key)
}
