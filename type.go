package request

import (
	"errors"
	"fmt"
	"net/http"
)

//go:generate gotemplate "./gotmpl/xchan" chanstr(string)
//go:generate gotemplate "./gotmpl/gomap" toreqs(string,MaybeReq)

type MaybeReq struct {
	Req *http.Request
	Err error
}

func ReqChan(urls []string, wrap func(string) MaybeReq) <-chan MaybeReq {
	return toreqs(1, wrap, chanstr(urls))
}

func GetRequests(urls []string) <-chan MaybeReq {
	return ReqChan(urls, DefaultGetRequest)
}

func DefaultGetRequest(url string) MaybeReq {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MaybeReq{nil, err}
	}
	req.Header.Add("Accept-Encoding", "gzip")
	return MaybeReq{req, nil}
}

type HttpData struct {
	Req *http.Request
	Res *http.Response
	Raw []byte
	Err error
}

type ResponseFilter func(*http.Response) error

func PassFilter() ResponseFilter {
	return func(res *http.Response) error {
		return nil
	}
}

func StatusFilter(codes []int) ResponseFilter {
	return func(res *http.Response) error {
		if res == nil {
			return errors.New("Response is nil")
		}
		for _, code := range codes {
			if code == res.StatusCode {
				return nil
			}
		}
		return fmt.Errorf("Expecting status code "+
			"in %v, found %d", codes, res.StatusCode)
	}
}
