package mux

import "net/http"

func MethodHandler(verb string, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == verb {
			h.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
