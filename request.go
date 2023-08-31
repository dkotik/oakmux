package oakmux

import "net/http"

// r.Body = http.MaxBytesReader(w, r.Body, MAX_API_PARAMS_SIZE)
// https://golang.hotexamples.com/examples/net.http/-/MaxBytesReader/golang-maxbytesreader-function-examples.html

type RequestReadLimiter struct {
	readLimit int64
	next      Handler
}

func NewRequestReadLimiter(next Handler, readLimit int64) *RequestReadLimiter {
	if next == nil {
		panic("cannot use a <nil> HTTP handler")
	}
	if readLimit == 0 {
		panic("cannot use a 0 read limit")
	}
	return &RequestReadLimiter{
		readLimit: readLimit,
		next:      next,
	}
}

func NewRequestReadLimiterMiddleware(readLimit int64) Middleware {
	return func(next Handler) Handler {
		return NewRequestReadLimiter(next, readLimit)
	}
}

func (l *RequestReadLimiter) ServeHyperText(
	w http.ResponseWriter,
	r *http.Request,
) error {
	r.Body = http.MaxBytesReader(w, r.Body, l.readLimit)
	return l.next.ServeHyperText(w, r)
}
