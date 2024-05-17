package url

type RequestBody struct {
	OriginalURL string `json:"original_url"`
}

type ResponseBody struct {
	OriginalURL  string `json:"original_url"`
	ShortenedUrl string `json:"short_url"`
}
