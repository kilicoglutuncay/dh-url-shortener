package service

import (
	"crypto/md5"
	"fmt"
)

type Shortener struct {
	Repository Repository
}

type Repository interface {
	Get(string) (string, error)
	Set(string, string) error
}

func (s Shortener) Shorten(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("long url cannot be empty")
	}

	shortURL := s.createShortURL(url, 0)
	return shortURL, nil
}

func (s Shortener) createShortURL(url string, collisionCounter int) string {
	input := []byte(url)
	counter := []byte(fmt.Sprintf("%d", collisionCounter))
	input = append(input, counter...)

	md5Sum := fmt.Sprintf("%x", md5.Sum(input))
	shortHash := md5Sum[:7]

	if err := s.Repository.Set(shortHash, url); err != nil {
		return s.createShortURL(url, collisionCounter+1)
	}

	return shortHash
}
