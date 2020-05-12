package emit

import (
	"sync"
	"testing"
)

func TestEmitter(t *testing.T) {
	e := New()

	var wg sync.WaitGroup

	fn1 := func(data interface{}) {
		t.Logf("fn1 receive event data: %v", data)
		wg.Done()
	}

	fn2 := func(data interface{}) {
		t.Logf("fn2 receive event data: %v", data)
		wg.Done()
	}

	fn3 := func(data interface{}) {
		t.Logf("fn3 receive event data: %v", data)
		wg.Done()
	}

	// test on
	e.On("testEvt", fn1)
	e.On("testEvt", fn2, fn3)
	wg.Add(3)

	// test emit
	e.Emit("testEvt", "this is data")
	wg.Wait()

	// test off
	e.Off("testEvt", fn3)
	wg.Add(2)
	e.Emit("testEvt", "this is data")
	wg.Wait()

}
