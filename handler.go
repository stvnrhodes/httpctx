// Package httpctx provides handlers to work with http and contexts.
package httpctx // import "stvn.cc/httpctx"

import (
	"net/http"

	"golang.org/x/net/context"
)

// Objects implementing the Handler interface can be
// registered to serve a particular path or subtree
// in the HTTP server.
//
// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return.  Returning signals that the request is finished
// and that the HTTP server can move on to the next request on
// the connection.
//
// If ServeHTTP panics, the server (the caller of ServeHTTP) assumes
// that the effect of the panic was isolated to the active request.
// It recovers the panic, logs a stack trace to the server error log,
// and hangs up the connection.
type Handler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	f(ctx, w, r)
}

// OldHandler converts from a http.Handler to a httpctx.Handler
func OldHandler(h http.Handler) Handler {
	return HandlerFunc(func(_ context.Context, w http.ResponseWriter, r *http.Request) { h.ServeHTTP(w, r) })
}

// OldHandleFunc converts from a http.HandlerFunc to a httpctx.HandlerFunc
func OldHandleFunc(f func(http.ResponseWriter, *http.Request)) HandlerFunc {
	return HandlerFunc(func(_ context.Context, w http.ResponseWriter, r *http.Request) { f(w, r) })
}

func rootHandler(h Handler) http.Handler {
	if h == nil {
		return nil
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		h.ServeHTTP(ctx, w, r)
	})
}

// ListenAndServe listens on the TCP network address addr and then calls Serve
// with handler to handle requests on incoming connections.
func ListenAndServe(addr string, handler Handler) error {
	return http.ListenAndServe(addr, rootHandler(handler))
}

// ListenAndServeTLS acts identically to ListenAndServe, except that it
// expects HTTPS connections. Additionally, files containing a certificate and
// matching private key for the server must be provided. If the certificate
// is signed by a certificate authority, the certFile should be the concatenation
// of the server's certificate, any intermediates, and the CA's certificate.
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler Handler) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, rootHandler(handler))
}
