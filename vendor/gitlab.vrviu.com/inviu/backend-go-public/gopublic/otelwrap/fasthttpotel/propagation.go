package fasthttpotel

import (
	"github.com/valyala/fasthttp"
)

func NewHeaderCarrier(h *fasthttp.RequestHeader) *HeaderCarrier {
	return &HeaderCarrier{header: h}
}

// HeaderCarrier adapts fasthttp.RequestHeader to satisfy the TextMapCarrier interface.
type HeaderCarrier struct {
	header *fasthttp.RequestHeader
}

// Get returns the value associated with the passed key.
func (hc *HeaderCarrier) Get(key string) string {
	return string(hc.header.Peek(key))
}

// Set stores the key-value pair.
func (hc *HeaderCarrier) Set(key string, value string) {
	hc.header.Set(key, value)
}

// Keys lists the keys stored in this carrier.
func (hc *HeaderCarrier) Keys() []string {
	keys := make([]string, 0, hc.header.Len())
	hc.header.VisitAll(func(key, val []byte) {
		keys = append(keys, string(key))
	})
	return keys
}
