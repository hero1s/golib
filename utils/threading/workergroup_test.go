package threading

import (
	"runtime"
	"sync"
	"testing"
)

func TestWorkerGroup(t *testing.T) {

	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())
	group := NewWorkerGroup(func() {
		lock.Lock()
		//do thing
		lock.Unlock()
		wg.Done()
	}, runtime.NumCPU())
	go group.Start()
	wg.Wait()

}
