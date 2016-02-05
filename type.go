package request

import "net/http"

type Request *http.Request

func GetRequests(urls []string) <-chan Request {
	out := make(chan Request)
	go func() {
		for _, url := range urls {
			if req, err := DefaultGetRequest(url); err == nil {
				out <- req
			}
		}
		close(out)
	}()
	return out
}

func DefaultGetRequest(url string) (Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept-Encoding", "gzip")
	return Request(req), err
}

type HttpData struct {
	Req *http.Request
	Res *http.Response
	Raw []byte
	Err error
}

type ResponseFilter func(*http.Response) bool

func PassFilter() ResponseFilter {
	return func(res *http.Response) bool {
		return true
	}
}

func StatusFilter(codes []int) ResponseFilter {
	return func(res *http.Response) bool {
		if res == nil {
			return false
		}
		for _, code := range codes {
			if code == res.StatusCode {
				return true
			}
		}
		return false
	}
}
