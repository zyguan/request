package request

import (
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
)

func Do(reqs <-chan Request, n int) <-chan HttpData {
	return DoWithFilter(reqs, n, PassFilter())
}

func DoWithFilter(reqs <-chan Request, n int,
	filter ResponseFilter) <-chan HttpData {

	out, n := make(chan HttpData), nGoroutines(n)
	if n == 1 {
		return Do1WithFilter(reqs, filter)
	}
	go func() {
		var wg sync.WaitGroup
		chans := make([]chan Request, n)
		cases := make([]reflect.SelectCase, n)

		// init and start workers
		for i := 0; i < n; i++ {
			chans[i] = make(chan Request)
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: reflect.ValueOf(chans[i]),
			}
			go fetchData(chans[i], out, &wg, filter)
		}

		// dispatch reqs to workers
		for req := range reqs {
			for i := range cases {
				cases[i].Send = reflect.ValueOf(req)
			}
			reflect.Select(cases)
		}

		// close channels
		for i := range chans {
			close(chans[i])
		}
		wg.Wait()
		close(out)
	}()
	return out
}

// Optimized version of Do(reqs, 1)
func Do1(reqs <-chan Request) <-chan HttpData {
	return Do1WithFilter(reqs, PassFilter())
}

func Do1WithFilter(reqs <-chan Request,
	filter ResponseFilter) <-chan HttpData {

	out := make(chan HttpData)
	go func() {
		fetchData(reqs, out, nil, filter)
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

	if wg != nil {
		wg.Add(1)
		defer wg.Done()
	}
	for req := range reqs {
		out <- RequestHttpData(req, filter)
	}
}

func RequestHttpData(req Request, filter ResponseFilter) HttpData {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return HttpData{req, res, []byte{}, err}
	}
	defer res.Body.Close()

	if !filter(res) {
		err = errors.New("Response doesn't match filter criteria")
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
