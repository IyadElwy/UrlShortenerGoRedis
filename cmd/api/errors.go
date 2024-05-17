package main

import (
	"net/http"
)

func returnErrorMessage(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}

func (app *application) notFoundError(w http.ResponseWriter) {
	returnErrorMessage(w, "Not found", http.StatusNotFound)
}

func (app *application) internalServerError(w http.ResponseWriter, err error) {
	app.errorLog.Println(err.Error())
	returnErrorMessage(w, "Internal server error", http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter) {
	returnErrorMessage(w, "Bad Request", http.StatusBadRequest)
}
