package utils

import "github.com/ReanSn0w/gokit/pkg/web"

type API interface {
	Request(string, string) *web.JsonRequest
}
