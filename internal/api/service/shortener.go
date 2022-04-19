package service

import (
	"crypto/sha256"
	"dh-url-shortener/internal/api/model"
	"fmt"
)

type Shortener struct {
	ShortURLDomain string
	DB             DB
}

type DB interface {
	Get(string) (model.RedirectionData, error)
	Set(string, model.RedirectionData) error
	Hit(string) error
	Data() map[string]model.RedirectionData
	Restore(map[string]model.RedirectionData)
}

// Shorten creates a short URL from a long URL
func (s Shortener) Shorten(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("long url cannot be empty")
	}

	hash := s.createShortURLHash(url, 0)
	shortURL := s.createShortURL(hash)
	return shortURL, nil
}

// createShortURLHash creates a hash from a long URL.
// Hash creation process is based on the following:
// 1. Create a SHA256 hash from the long URL with collision counter
// 2. Pick first seven character of the hash as the short URL
// 3. If the short URL is already taken, create a new hash with collision counter and repeat the process

func (s Shortener) createShortURLHash(url string, collisionCounter int) string {
	input := []byte(url)
	counter := []byte(fmt.Sprintf("%d", collisionCounter))
	input = append(input, counter...)

	hash := fmt.Sprintf("%x", sha256.Sum256(input))
	shortHash := hash[:7]

	if err := s.DB.Set(shortHash, model.RedirectionData{OriginalURL: url, Hits: 0}); err != nil {
		return s.createShortURLHash(url, collisionCounter+1)
	}

	return shortHash
}

// createShortURL creates a short URL from short URL domain and hash
func (s Shortener) createShortURL(hash string) string {
	return fmt.Sprintf("%s/%s", s.ShortURLDomain, hash)
}

// Expand expands a short URL to a long URL
func (s Shortener) Expand(hash string) (string, error) {
	redirectionData, err := s.DB.Get(hash)
	if err != nil {
		return "", err
	}

	err = s.DB.Hit(hash)
	if err != nil {
		return "", err
	}

	return redirectionData.OriginalURL, nil
}

// List converts the data in the database to a list of data which contains short URL, long URL and hit count
func (s Shortener) List() []model.ListData {
	data := s.DB.Data()
	list := make([]model.ListData, 0, len(data))
	for k, v := range data {
		list = append(list, model.ListData{Hash: k, OriginalURL: v.OriginalURL, Hits: v.Hits})
	}
	return list
}
