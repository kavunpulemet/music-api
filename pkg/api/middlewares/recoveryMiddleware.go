package middlewares

import (
	"net/http"
	"runtime/debug"

	"effectiveMobileTest/pkg/api/utils"
)

func RecoveryMiddleware(ctx utils.MyContext, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ctx.Logger.Errorf("panic occurred: %v\n%s", err, debug.Stack())
				utils.NewErrorResponse(ctx, w, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}
