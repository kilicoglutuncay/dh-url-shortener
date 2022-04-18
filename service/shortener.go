package service

import (
	"crypto/sha256"
	"fmt"
)

type Shortener struct {
	ShortURLDomain string
	DB             DB
}

type DB interface {
	Get(string) (string, error)
	Set(string, string) error
	Data() map[string]string
	Restore(map[string]string)
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

	if err := s.DB.Set(shortHash, url); err != nil {
		return s.createShortURLHash(url, collisionCounter+1)
	}

	return shortHash
}

func (s Shortener) createShortURL(hash string) string {
	return fmt.Sprintf("%s/%s", s.ShortURLDomain, hash)
}

func (s Shortener) Expand(hash string) (string, error) {
	url, err := s.DB.Get(hash)
	if err != nil {
		return "", err
	}

	return url, nil
}
