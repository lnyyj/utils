package zksrv

import (
	"github.com/samuel/go-zookeeper/zk"
	"github.com/sirupsen/logrus"
)

// ConfCenter 配置中心
type ConfCenter struct {
	zkSrv

	flag int8 // 0 监听当前节点， 1 监听当前节点的子节点(仅创建可检测)
}

// NewCC 创建一个配置中心对象
func NewCC(hosts ...string) *ConfCenter {
	cc := &ConfCenter{}
	cc.SetHosts(hosts...).CreateZKCli()
	return cc
}

// SetCurrentWFlag 设置监听当前节点
func (cc *ConfCenter) SetCurrentWFlag() *ConfCenter {
	cc.flag = 0
	return cc
}

// SetChildrenWFlag 设置监听当前节点的子节点(仅创建可检测)
func (cc *ConfCenter) SetChildrenWFlag() *ConfCenter {
	cc.flag = 1
	return cc
}

// CC 启动配置中心
func (cc *ConfCenter) CC(path string) (data chan string, errs chan error) {
	return cc.listenWatch(path)
}

// ListenWatch 监听一个路径的watch事件
func (cc *ConfCenter) listenWatch(path string) (chan string, chan error) {
	cc.checkzk()

	data := make(chan string)
	errs := make(chan error)

	logrus.Info("listen watch start")
	go func() {
		var err error
		var events <-chan zk.Event
		for {
			select {
			case <-cc.stop:
				logrus.Info("listen watch stop")
				return
			default:
				if cc.flag == 0 {
					_, _, events, err = cc.zkConn.ExistsW(path)
				} else {
					_, _, events, err = cc.zkConn.ChildrenW(path)
				}

				if err != nil {
					errs <- err
					return
				}

				evt := <-events
				logrus.Debug("event:[%v]", evt)

				if evt.Err != nil {
					errs <- evt.Err
					return
				}

				go func(event zk.Event) {
					data <- cc.watchCB(event)
				}(evt)
			}
		}
	}()

	return data, errs
}

// WatchCB watch事件回调/处理函数
func (cc *ConfCenter) watchCB(event zk.Event) string {
	logrus.Debug("callback event:[%v]", event)

	// 节点数据和子节点数据更新才重新获取
	if event.Type != zk.EventNodeDataChanged &&
		event.Type != zk.EventNodeChildrenChanged {
		logrus.Info("invalid event type:[%s]", event.Type.String())
		return ""
	}

	// 更新全局数据
	data, stat, err := cc.zkConn.Get(event.Path)
	if err != nil {
		logrus.Error("Get path:[%s] returned error:[%s]", event.Path, err.Error())
		return ""
	} else if stat == nil {
		logrus.Error("Get path:[%s] returned nil stat", event.Path)
		return ""
	} else if len(data) < 4 {
		logrus.Error("Get path:[%s] returned wrong size data", event.Path)
		return ""
	}

	return string(data)
}
