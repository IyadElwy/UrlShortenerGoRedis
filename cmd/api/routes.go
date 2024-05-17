package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.GET("/*urlId", app.resolveUrl)
	router.POST("/create", app.createNewShortUrl)
	return app.logRequestMiddleware(router)
}
