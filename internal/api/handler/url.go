package handler

import (
	"dh-url-shortener/internal/api/model"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type URLHandler struct {
	ShortenerService ShortenerService
}

type ShortenerService interface {
	Shorten(string) (string, error)
	Expand(string) (string, error)
	List() []model.ListData
}

const (
	shortURLHashLength = 7
	errInvalidURL      = "invalid url"
)

// Shorten handles requests which are aim to shorten long URL.
func (h URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var sr ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = sr.validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL, err := h.ShortenerService.Shorten(sr.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.json(w, http.StatusCreated, &ShortenResponse{URL: shortURL})
}

// Expand expands the given short URL to its long URL.
func (h URLHandler) Expand(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[1:]
	if len(hash) != shortURLHashLength {
		http.Error(w, errors.New("invalid hash").Error(), http.StatusBadRequest)
		return
	}

	longURL, err := h.ShortenerService.Expand(hash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

// List returns a list of all stored URLs with their hits.
func (h URLHandler) List(w http.ResponseWriter, _ *http.Request) {
	listData := h.ShortenerService.List()
	h.json(w, http.StatusOK, &listData)
}

// json encodes the given data to JSON and writes it to the response writer.
func (h URLHandler) json(w http.ResponseWriter, statusCode int, data interface{}) {
	resp, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(resp)
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	URL string `json:"url"`
}

// Validate validates the ShortenRequest.URL field is a valid URL
func (r ShortenRequest) validate() error {
	if r.URL == "" {
		return errors.New(errInvalidURL)
	} else if _, err := url.ParseRequestURI(r.URL); err != nil {
		return err
	}

	return nil
}
