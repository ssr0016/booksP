package response

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

var (
	_ http.ResponseWriter = &responseWriter{}
	_ http.Hijacker       = &responseWriter{}
)

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	// Status returns the status code of the response or 0 if the response has not been written.
	Status() int
	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
	// Size returns the size of the response body.
	Size() int

	Error(err error) error
}

// BeforeFunc is a function that is called before the ResponseWriter has been written to.
type BeforeFunc func(ResponseWriter)

// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewResponseWriter(method string, rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{method, rw, 0, 0, nil, nil}
}

// Rw returns a ResponseWriter. If the argument already satisfies the interface,
// it is returned as is, otherwise it is wrapped using NewResponseWriter
func Rw(rw http.ResponseWriter, req *http.Request) ResponseWriter {
	if mrw, ok := rw.(ResponseWriter); ok {
		return mrw
	}

	return NewResponseWriter(req.Method, rw)
}

type responseWriter struct {
	method string
	http.ResponseWriter
	status      int
	size        int
	beforeFuncs []BeforeFunc
	err         error
}

func (rw *responseWriter) WriteHeader(s int) {
	rw.callBefore()

	// Avoid panic if status code is not a valid HTTP status code
	if s < 100 || s > 999 {
		rw.ResponseWriter.WriteHeader(500)
		rw.status = 500
		return
	}

	rw.ResponseWriter.WriteHeader(s)
	rw.status = s
}

func (rw *responseWriter) Write(b []byte) (size int, err error) {
	if !rw.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.WriteHeader(http.StatusOK)
	}
	if rw.method != "HEAD" {
		size, err = rw.ResponseWriter.Write(b)
		rw.size += size
	}
	return size, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) Error(err error) error {
	if err != nil {
		rw.err = err
	}
	return rw.err
}

func (rw *responseWriter) Before(before BeforeFunc) {
	rw.beforeFuncs = append(rw.beforeFuncs, before)
}

const StatusHijacked = -1

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("the ResponseWriter doesn't support the Hijacker interface")
	}

	conn, brw, err := hijacker.Hijack()
	if err == nil {
		rw.status = StatusHijacked
	}

	return conn, brw, err
}

func (rw *responseWriter) callBefore() {
	for i := len(rw.beforeFuncs) - 1; i >= 0; i-- {
		rw.beforeFuncs[i](rw)
	}
}

func (rw *responseWriter) Flush() {
	flusher, ok := rw.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}
