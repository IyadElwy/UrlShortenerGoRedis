package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"urlShortener.IyadElwy/internal/db"
	"urlShortener.IyadElwy/internal/url"
)

func (app *application) createNewShortUrl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			app.clientError(w)
			return
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) //1mb

	shortId, err := app.idGenerator.Generate()
	if err != nil {
		app.internalServerError(w, err)
		return
	}
	var requestBody url.RequestBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = app.redis.Store(shortId, requestBody.OriginalURL)
	if err != nil {
		app.internalServerError(w, err)
		return
	}

	protocol, ok := app.envs["protocol"]
	if !ok {
		app.internalServerError(w, fmt.Errorf("Protocol not set correctly in .env"))
	}
	responseBody := url.ResponseBody{
		OriginalURL:  requestBody.OriginalURL,
		ShortenedUrl: fmt.Sprintf("%s://%s/%s", protocol, r.Host, shortId),
	}
	jsonResponse, err := json.MarshalIndent(&responseBody, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonResponse))
}

func (app *application) resolveUrl(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	urlID := strings.Split(params.ByName("urlId"), "/")[1]
	resolvedUrl, err := app.redis.Retrieve(urlID)
	if err != nil {
		switch {
		case err == db.ErrNotFoundInRedis:
			app.notFoundError(w)
		default:
			app.internalServerError(w, err)
		}
		return
	}
	http.Redirect(w, r, resolvedUrl, http.StatusSeeOther)
}
