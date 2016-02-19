package request

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
)

var urls = []string{
	"https://golang.org/",
	"https://kotlinlang.org/",
	"https://www.rust-lang.org/",
	"https://www.haskell.org/",
	"https://racket-lang.org/",
	"https://www.python.org/",
	"https://www.ruby-lang.org/",
	"https://github.com/zyguan/None",
}
var pad = padding(urls)

func padding(strs []string) func(string) string {
	l := 0
	for _, s := range strs {
		if len(s) > l {
			l = len(s)
		}
	}
	format := "%-" + strconv.Itoa(l) + "s"
	return func(s string) string {
		return fmt.Sprintf(format, s)
	}
}

func TestDo(t *testing.T) {
	recv(Do(GetRequests(urls), runtime.NumCPU()))
}

func TestDoWithFilter(t *testing.T) {
	recv(DoWithFilter(GetRequests(urls), 1, StatusFilter([]int{200})))
}

func recv(data <-chan HttpData) {
	fmt.Printf("recv from %#v\n", data)
	for d := range data {
		if d.Err != nil {
			fmt.Println(pad(d.Req.URL.String()), d.Err)
		} else {
			fmt.Println(pad(d.Req.URL.String()), d.Res.Status)
		}
	}
}
