package request

import "sync"

// template type GoMap(A, B)

func doRequest(n int, f func(MaybeReq) HttpData, in <-chan MaybeReq) <-chan HttpData {
	if n <= 0 {
		n = 1
	}
	out := make(chan HttpData)
	go func() {
		var wg sync.WaitGroup
		wg.Add(n)
		for i := 0; i < n; i++ {
			go func() {
				for a := range in {
					out <- f(a)
				}
				wg.Done()
			}()
		}
		wg.Wait()
		close(out)
	}()
	return out
}
