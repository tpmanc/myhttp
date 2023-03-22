package http

import "net/url"

type Mock struct {
	DoRequestFn func(u *url.URL) ([]byte, error)
}

func (m Mock) DoRequest(u *url.URL) ([]byte, error) {
	return m.DoRequestFn(u)
}
