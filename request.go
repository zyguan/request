package request

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

func Do(reqs <-chan Request, n int) <-chan HttpData {
	return DoWithFilter(reqs, n, PassFilter())
}

func DoWithFilter(reqs <-chan Request, n int,
	filter ResponseFilter) <-chan HttpData {

	out, n := make(chan HttpData), nGoroutines(n)

	go func() {
		var wg sync.WaitGroup
		wg.Add(n)
		for i := 0; i < n; i++ {
			go fetchData(reqs, out, &wg, filter)
		}
		wg.Wait()
		close(out)
	}()
	return out
}

func nGoroutines(n int) int {
	if n <= 0 {
		return 1
	}
	return n
}

func fetchData(reqs <-chan Request, out chan<- HttpData,
	wg *sync.WaitGroup, filter ResponseFilter) {
	for req := range reqs {
		out <- RequestHttpData(req, filter)
	}
	wg.Done()
}

func RequestHttpData(req Request, filter ResponseFilter) HttpData {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return HttpData{req, res, []byte{}, err}
	}
	defer res.Body.Close()

	if err := filter(res); err != nil {
		return HttpData{req, res, []byte{}, err}
	}

	in, err := decodeBody(res)
	if err != nil {
		return HttpData{req, res, []byte{}, err}
	}

	raw, err := ioutil.ReadAll(in)
	return HttpData{req, res, raw, err}
}

func decodeBody(res *http.Response) (io.Reader, error) {
	switch enc := res.Header.Get("Content-Encoding"); enc {
	case "gzip":
		return gzip.NewReader(res.Body)
	default:
		return res.Body, nil
	}
}
