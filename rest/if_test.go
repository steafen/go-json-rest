package rest

import (
	"testing"

	"github.com/AlexanderChen1989/go-json-rest/rest/test"
	"golang.org/x/net/context"
)

func TestIfMiddleware(t *testing.T) {

	api := NewApi()

	// the middleware to test
	api.Use(&IfMiddleware{
		Condition: func(r *Request) bool {
			if r.URL.Path == "/true" {
				return true
			}
			return false
		},
		IfTrue: MiddlewareSimple(func(handler HandlerFunc) HandlerFunc {
			return func(ctx context.Context, w ResponseWriter, r *Request) {
				r.Env["TRUE_MIDDLEWARE"] = true
				handler(ctx, w, r)
			}
		}),
		IfFalse: MiddlewareSimple(func(handler HandlerFunc) HandlerFunc {
			return func(ctx context.Context, w ResponseWriter, r *Request) {
				r.Env["FALSE_MIDDLEWARE"] = true
				handler(ctx, w, r)
			}
		}),
	})

	// a simple app
	api.SetApp(AppSimple(func(ctx context.Context, w ResponseWriter, r *Request) {
		w.WriteJson(r.Env)
	}))

	// wrap all
	handler := api.MakeHandler()

	recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs("{\"FALSE_MIDDLEWARE\":true}")

	recorded = test.RunRequest(t, handler, test.MakeSimpleRequest("GET", "http://localhost/true", nil))
	recorded.CodeIs(200)
	recorded.ContentTypeIsJson()
	recorded.BodyIs("{\"TRUE_MIDDLEWARE\":true}")
}
