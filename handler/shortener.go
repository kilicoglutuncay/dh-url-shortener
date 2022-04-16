package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type ShortenerHandler struct {
	ShortenerService ShortenerService
}

type ShortenerService interface {
	Shorten(string) (string, error)
}

// Create creates a new short URL
func (h ShortenerHandler) Create(w http.ResponseWriter, r *http.Request) {

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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(shortURL))
}

const ErrInvalidURL = "invalid url"

type ShortenRequest struct {
	URL string `json:"url"`
}

// Validate validates the ShortenRequest.URL field is a valid URL
func (r ShortenRequest) validate() error {
	if r.URL == "" {
		return errors.New(ErrInvalidURL)
	} else if _, err := url.ParseRequestURI(r.URL); err != nil {
		return err
	}

	return nil
}
