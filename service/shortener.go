package service

import (
	"crypto/sha256"
	"dh-url-shortener/model"
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

func (s Shortener) Shorten(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("long url cannot be empty")
	}

	hash := s.createShortURLHash(url, 0)
	shortURL := s.createShortURL(hash)
	return shortURL, nil
}

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

func (s Shortener) createShortURL(hash string) string {
	return fmt.Sprintf("%s/%s", s.ShortURLDomain, hash)
}

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

func (s Shortener) List() []model.ListData {
	data := s.DB.Data()
	var list []model.ListData
	for k, v := range data {
		list = append(list, model.ListData{Hash: k, OriginalURL: v.OriginalURL, Hits: v.Hits})
	}
	return list
}
