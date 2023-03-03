package middleware

import (
	"net/http"
	"walltrack/util"
)

func JsonRoute(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			util.WriteApiErrMessage(w, 0, "Invalid content type")
		}

		h.ServeHTTP(w, r)
	})
}
