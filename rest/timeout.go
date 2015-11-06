package rest

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

// TimeoutMiddleware cancel context when timeout
type TimeoutMiddleware struct {
	timeout time.Duration
}

// Timeout create timeout middleware with duration
func Timeout(timeout time.Duration) *TimeoutMiddleware {
	return &TimeoutMiddleware{timeout}
}

// MiddlewareFunc makes TimeoutMiddleware implement the Middleware interface.
func (mw *TimeoutMiddleware) MiddlewareFunc(h HandlerFunc) HandlerFunc {
	return func(ctx context.Context, w ResponseWriter, r *Request) {
		ctx, _ = context.WithTimeout(ctx, mw.timeout)
		h(ctx, w, r)
		// Cancel the context if the client closes the connection
		if wcn, ok := w.(http.CloseNotifier); ok {
			var cancel context.CancelFunc
			ctx, cancel = context.WithCancel(ctx)
			defer cancel()

			notify := wcn.CloseNotify()
			go func() {
				<-notify
				cancel()
			}()
		}

		h(ctx, w, r)
	}
}
