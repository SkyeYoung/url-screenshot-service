package helper

import (
	"errors"

	"strings"
)

type fnKV struct {
	k func(url string) bool
	v func()
}

func GetValidUrl(urlStr string) (u string, err error) {

	arr := []fnKV{
		{k: func(u string) bool { return strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") }, v: func() { u = urlStr }},
		{k: func(u string) bool { return strings.HasPrefix(u, "//") }, v: func() { u = "http:" + urlStr }},
		{k: func(u string) bool { return strings.Contains(u, "//") }, v: func() { u, err = "", errors.New("invalid http/https url: "+urlStr) }},
		{k: func(u string) bool { return true }, v: func() { u = "http://" + urlStr }},
	}

	for _, kv := range arr {
		if kv.k(urlStr) {
			kv.v()
			break
		}
	}

	return
}
