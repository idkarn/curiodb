package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RouteHandler func(RequestContext)

type Route struct {
	Method  string
	Path    string
	Handler RouteHandler
}

func NewRouteInfo(method, path string, handler RouteHandler) Route {
	return Route{method, path, handler}
}

type RequestContext struct {
	Route    Route
	Request  *http.Request
	Response http.ResponseWriter
	Data     any
	phase    uint8 // 0 - middleware; 1 - main handler
}

func NewRequestContext(route Route, req *http.Request, res http.ResponseWriter) RequestContext {
	return RequestContext{route, req, res, nil, 0}
}

func (ctx *RequestContext) Send(resp any) {
	fmt.Fprint(ctx.Response, resp)
}

func (ctx *RequestContext) SendBytes(bytes []byte) {
	ctx.Response.Write(bytes)
}

func (ctx *RequestContext) SendJSON(response any) {
	out, err := json.Marshal(response)
	if err != nil {
		panic("Unable to Marshal this object")
	}
	ctx.SendBytes(out)
}

func (ctx *RequestContext) Status(statusCode int) {
	ctx.Response.WriteHeader(statusCode)
}

func (ctx *RequestContext) Error(statusMessage string, statusCode int) {
	http.Error(ctx.Response, statusMessage, statusCode)
}

func (ctx *RequestContext) Read(dest any) error {
	err := json.NewDecoder(ctx.Request.Body).Decode(dest)
	if err != nil {
		return errors.New("unable to decode this json")
	}
	ctx.Data = dest
	return nil
}

type NextFunction func()

type MiddlewareFn func(RequestContext, NextFunction)

var Middlewares []MiddlewareFn

func newNextFunction(ctx *RequestContext, fnIdx int) NextFunction {
	if fnIdx != len(Middlewares) {
		return func() {
			Middlewares[fnIdx](*ctx, newNextFunction(ctx, fnIdx+1))
		}
	} else {
		return func() {
			ctx.phase = 1
		}
	}
}

func RunWith(ctx RequestContext) bool {
	if len(Middlewares) == 0 {
		return true
	}
	next := newNextFunction(&ctx, 0)
	next()
	return ctx.phase == 1
}

func HandleWith(w http.ResponseWriter, r *http.Request, route Route) {
	ctx := NewRequestContext(route, r, w)
	if ok := RunWith(ctx); ok {
		route.Handler(ctx)
	}
}

func SetupMiddlewares(config []MiddlewareFn) {
	Middlewares = append(Middlewares, config...)
}
