package repository

import "errors"

type InMemoryRepository struct {
	data map[string]string
}

func NewInMemoryRepository(data map[string]string) *InMemoryRepository {
	return &InMemoryRepository{
		data: data,
	}
}

func (i InMemoryRepository) Get(s string) (string, error) {
	return "", nil
}

func (i InMemoryRepository) Set(key string, value string) error {
	if _, ok := i.data[key]; ok {
		return errors.New(key + " key already exists")
	}
	i.data[key] = value
	return nil
}
