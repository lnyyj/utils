package works

import (
	"fmt"
	"testing"
	"time"
)

type WI int

func (wi *WI) Process() (IWork, error) {
	fmt.Printf("--->wi:[%d]\r\n", *wi)
	return WO(*wi), nil
}

type WO int

func (wo WO) Process() (IWork, error) {
	fmt.Printf("--->wo:[%d]\r\n", wo)
	return nil, nil
}

func Test_Works(t *testing.T) {
	inqueue := NewChannel(30)
	outqueue := NewChannel(30)
	work1 := NewDispatcher(2, inqueue, outqueue).Run()
	work2 := NewDispatcher(2, outqueue, nil).Run()

	for i := 0; i < 5; i++ {
		wi := WI(i)
		inqueue <- &wi
		time.Sleep(100 * time.Millisecond)
	}

	work1.Close()
	work2.Close()
}
