package mux

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/ghost/handlers"
)

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

const StatusUpgradeRequired = 426

func ErrorHandler(w http.ResponseWriter, err error) {
	log.Printf("Error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
}

func HttpsRequiredHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if !(r.URL.Scheme == "https" || r.Header.Get("X-Forwarded-Proto") == "https") {
			log.Printf("Invalid scheme: %v/%v", r.URL, r.Header.Get("X-Forwarded-Proto"))
			w.WriteHeader(StatusUpgradeRequired)
		} else {
			h.ServeHTTP(w, r)
		}
	}

	return http.HandlerFunc(fn)
}

type AuthFunc func(string, string) bool

func BasicAuthHandler(authFunc AuthFunc, realm string, h http.Handler) http.Handler {
	authFn := func(user, pass string) (interface{}, bool) {
		return nil, authFunc(user, pass)
	}

	return handlers.BasicAuthHandler(h, authFn, realm)
}
