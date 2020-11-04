package kvika

import (
	"net/http"
	"net/url"
)

type Request struct {
	Method  string
	URL     *url.URL
	Headers http.Header
	Payload []byte
}
