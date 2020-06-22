package cookie

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

var (
	// Interface guards to ensure this module implements the following interfaces
	_ caddy.Provisioner           = (*Cookie)(nil)
	_ caddy.Validator             = (*Cookie)(nil)
	_ caddyhttp.MiddlewareHandler = (*Cookie)(nil)
)

func init() {
	caddy.RegisterModule(Cookie{})
}

type Cookie struct {
	logger          *zap.Logger
	Disabled		bool
	CookiesToRemove []string
}

func (c Cookie) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.cookie",
		New: func() caddy.Module { return new(Cookie) },
	}
}

// Provision sets up the module.
func (c *Cookie) Provision(ctx caddy.Context) error {
	c.logger = ctx.Logger(c) // g.logger is a *zap.Logger
	return nil
}

// Validate validates that the module has a usable config.
func (c Cookie) Validate() error {
	// TODO: validate the module's setup
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (c Cookie) ServeHTTP(writer http.ResponseWriter, request *http.Request, nextHandler caddyhttp.Handler) error {
	if !c.Disabled {
		if len(c.CookiesToRemove) > 0 {
			removeCookiesFromRequest(c.CookiesToRemove, request)
		}
	}
	return nextHandler.ServeHTTP(writer, request)
}

func removeCookiesFromRequest(cookiesToRemove []string, request *http.Request) {
	cookies := request.Cookies()
	if len(cookies) == 0 {
		return
	}
	if len(cookiesToRemove) == 0 {
		return
	}
	cookiesToKeep := cookies[:0]
	for _, cookie := range cookies {
		_, found := find(cookiesToRemove, cookie.Name)
		if !found {
			cookiesToKeep = append(cookiesToKeep, cookie)
		}
	}
	strs := make([]string, len(cookiesToKeep))
	for i, v := range cookiesToKeep {
		strs[i] = v.String()
	}
	request.Header.Set("Cookie", strings.Join(strs, ";"))
}

// linear search
func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if strings.EqualFold(item, val) { // case-insensitive
			return i, true
		}
	}
	return -1, false
}