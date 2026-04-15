package xsync

import (
	"sync"
)

func WaitAll(fns ...func()) {
	wg := &sync.WaitGroup{}
	wg.Add(len(fns))
	for _, fn := range fns {
		go func() {
			defer wg.Done()
			fn()
		}()
	}
	wg.Wait()
}
