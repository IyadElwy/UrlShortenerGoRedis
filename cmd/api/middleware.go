package main

import (
	"net/http"

	"github.com/urfave/negroni"
)

func (app *application) logRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := negroni.NewResponseWriter(w)
		next.ServeHTTP(lrw, r)
		app.infoLog.Printf("%s|%s|%s|%d", r.RemoteAddr, r.Method, r.RequestURI, lrw.Status())
	})
}
