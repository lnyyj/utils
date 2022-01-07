package laterinit

import (
	"log"
)

type initFunc func()

type initModule struct {
	name string
	fn   initFunc
}

var initFuncs []*initModule

func init() {
	initFuncs = make([]*initModule, 0)
}

func RegisterInit(name string, fn initFunc) {
	log.Println("register module ", name)
	initFuncs = append(initFuncs, &initModule{
		name: name,
		fn:   fn,
	})
}

func Init() {
	for _, service := range initFuncs {
		log.Println("init module  ", service.name)
		service.fn()
	}
}
