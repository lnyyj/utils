package laterinit

import (
	"log"
)


// init module
type initModule struct {
	name string
	fn   func() error
}

var initFuncs []*initModule

func init() {
	initFuncs = make([]*initModule, 0)
}

func Register(name string, fn func() error) {
	log.Println("register module ", name)
	initFuncs = append(initFuncs, &initModule{
		name: name,
		fn:   fn,
	})
}

func Init() {
	for _, service := range initFuncs {
		log.Println("init module  ", service.name)
		if err := service.fn(); err != nil {
			panic(err)
		}
	}
}
