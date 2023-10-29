package middleware

import "net/http"

func CheckRouteMethod(ctx RequestContext, next NextFunction) {
	if ctx.Request.Method != ctx.Route.Method {
		ctx.Error("Unsupported method", http.StatusMethodNotAllowed)
	} else {
		next()
	}
}
