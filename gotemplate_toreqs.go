package request

import "sync"

// template type GoMap(A, B)

func toreqs(n int, f func(string) MaybeReq, in <-chan string) <-chan MaybeReq {
	if n <= 0 {
		n = 1
	}
	out := make(chan MaybeReq)
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
