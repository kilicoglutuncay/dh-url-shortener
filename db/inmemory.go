package db

import (
	"errors"
)

type InMemoryDB struct {
	data map[string]string
}

func NewInMemoryDB() *InMemoryDB {
	repo := &InMemoryDB{
		data: make(map[string]string),
	}

	return repo
}

func (i InMemoryDB) Get(key string) (string, error) {
	if value, ok := i.data[key]; ok {
		return value, nil
	}
	return "", errors.New(key + " not found")
}

func (i InMemoryDB) Set(key, value string) error {
	if _, ok := i.data[key]; ok {
		return errors.New(key + " key already exists")
	}
	i.data[key] = value
	return nil
}

func (i InMemoryDB) Data() map[string]string {
	return i.data
}

func (i *InMemoryDB) Restore(data map[string]string) {
	i.data = data
}
