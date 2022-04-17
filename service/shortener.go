package service

import (
	"crypto/md5"
	"fmt"
)

type Shortener struct {
	ShortURLDomain string
	Repository     Repository
}

type Repository interface {
	Get(string) (string, error)
	Set(string, string) error
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

	md5Sum := fmt.Sprintf("%x", md5.Sum(input))
	shortHash := md5Sum[:7]

	if err := s.Repository.Set(shortHash, url); err != nil {
		return s.createShortURLHash(url, collisionCounter+1)
	}

	return shortHash
}

func (s Shortener) createShortURL(hash string) string {
	return fmt.Sprintf("%s/%s", s.ShortURLDomain, hash)
}

func (s Shortener) Expand(shortURL string) (string, error) {
	return "", nil
}
