package request

// template type XChan(A)

func chanstr(xs []string) <-chan string {
	out := make(chan string)
	go func() {
		for _, x := range xs {
			out <- x
		}
		close(out)
	}()
	return out
}
