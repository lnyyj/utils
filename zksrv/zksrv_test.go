package zksrv

import (
	"testing"

	"github.com/sirupsen/logrus"
)

// Test_CC 配置中心测试
func Test_CC(t *testing.T) {
	cc := NewCC("192.168.1.15:2181")
	// cc.SetHosts("192.168.1.115:2181").CreateZKCli()

	go func() {
		d, err := cc.CC("/conf")
		for {
			for {
				select {
				case d := <-d:
					logrus.Debug("--->data: %+v\n", d)
				case err := <-err:
					logrus.Debug("--->err: %+v\n", err)
				}
			}
		}
	}()

	// go func() {
	// 	d, err := cc.SetChildrenWFlag().CC("/conf")
	// 	for {
	// 		for {
	// 			select {
	// 			case d := <-d:
	// 				logrus.Debug("--->children data: %+v\n", d)
	// 			case err := <-err:
	// 				logrus.Debug("--->children err: %+v\n", err)
	// 			}
	// 		}
	// 	}
	// }()

	a := make(chan int)
	<-a
}
