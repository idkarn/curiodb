package middleware

import "net/http"

type APIRequest struct {
	Request  *http.Request
	Response http.ResponseWriter
}
type MiddlewareFn func(APIRequest)

type DecodedRequest struct {
	Response    http.Response
	DecodedData interface{}
}

var Middlewares []MiddlewareFn
var mwIdx int = 0

func RegisterMiddleware(fn MiddlewareFn) {
	Middlewares = append(Middlewares, fn)
}

func next(r APIRequest) {
	// !ATTENTION: concurrent execution will be able to overwrite counter value
	if mwIdx == len(Middlewares) {
		return
	}
	Middlewares[mwIdx](r)
	mwIdx++
}

func ExecuteMuddlewares(r APIRequest) {
	if len(Middlewares) == 0 {
		return
	}
	next(r)
}

func SetupMiddlewares(config []MiddlewareFn) {
	for _, fn := range config {
		RegisterMiddleware(fn)
	}
}

func DecodeRequestPayload(r APIRequest) {

}
