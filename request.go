package request

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
)

//go:generate gotemplate "./gotmpl/gomap" doRequest(MaybeReq,HttpData)

func Do(n int, reqs <-chan MaybeReq) <-chan HttpData {
	return DoWithFilter(n, reqs, PassFilter())
}

func DoWithFilter(n int, reqs <-chan MaybeReq,
	filter ResponseFilter) <-chan HttpData {
	return doRequest(n, HttpRequester(filter), reqs)
}

func HttpRequester(filter ResponseFilter) func(MaybeReq) HttpData {
	return func(r MaybeReq) HttpData {
		req, err := r.Req, r.Err
		if err != nil {
			return HttpData{req, nil, nil, err}
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return HttpData{req, res, nil, err}
		}
		defer res.Body.Close()

		if err := filter(res); err != nil {
			return HttpData{req, res, nil, err}
		}

		in, err := decodeBody(res)
		if err != nil {
			return HttpData{req, res, nil, err}
		}

		raw, err := ioutil.ReadAll(in)
		return HttpData{req, res, raw, err}
	}
}

func decodeBody(res *http.Response) (io.Reader, error) {
	switch enc := res.Header.Get("Content-Encoding"); enc {
	case "gzip":
		return gzip.NewReader(res.Body)
	default:
		return res.Body, nil
	}
}
