package cookie

import (
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware_ServeHTTP(t *testing.T) {
	type fields struct {
		logger *zap.Logger
		disabled bool
		cookiesToRemove []string
	}
	type args struct {
		writer      http.ResponseWriter
		request     *http.Request
		nextHandler caddyhttp.Handler
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Given cookies in request that should be removed, when ServeHTTP, then cookies are removed from request.",
			fields{
				logger: zaptest.NewLogger(t),
				cookiesToRemove: []string{"cookie-to-remove"},
			},
			args{
				writer: httptest.NewRecorder(),
				request: &http.Request{
					Header: http.Header{
						"Cookie": []string{"cookie-to-remove=cookie-value","cookie=cookie-value; cook2=value"},
					},
				},
				nextHandler: AssertionHandler{
					t:    t,
					assertions: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
						_, err := r.Cookie("cookie-to-remove")
						assert.Error(t, err)
						cookie, _ := r.Cookie("cookie")
						assert.Equal(t, "cookie-value", cookie.Value)
						cookie, _ = r.Cookie("cook2")
						assert.Equal(t, "value", cookie.Value)
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := Cookie{
				logger: tt.fields.logger,
				Disabled: tt.fields.disabled,
				CookiesToRemove: tt.fields.cookiesToRemove,
			}
			if err := x.ServeHTTP(tt.args.writer, tt.args.request, tt.args.nextHandler); (err != nil) != tt.wantErr {
				t.Errorf("ServeHTTP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type AssertionHandler struct {
	t    *testing.T
	assertions assertFunc
}

type assertFunc func(t *testing.T, w http.ResponseWriter, r *http.Request)

func (a AssertionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) error {
	a.assertions(a.t, writer, request)
	return nil
}