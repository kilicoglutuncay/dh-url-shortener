package db

import (
	"dh-url-shortener/model"
	"errors"
)

type InMemoryDB struct {
	data map[string]model.RedirectionData
}

func NewInMemoryDB() *InMemoryDB {
	repo := &InMemoryDB{
		data: make(map[string]model.RedirectionData),
	}

	return repo
}

func (i InMemoryDB) Get(key string) (model.RedirectionData, error) {
	if value, ok := i.data[key]; ok {
		return value, nil
	}
	return model.RedirectionData{}, errors.New(key + " not found")
}

func (i InMemoryDB) Set(key string, value model.RedirectionData) error {
	if _, ok := i.data[key]; ok {
		return errors.New(key + " key already exists")
	}
	i.data[key] = value
	return nil
}

func (i InMemoryDB) Hit(key string) error {
	value, ok := i.data[key]
	if !ok {
		return errors.New(key + " key not exists")
	}
	value.Hits++
	i.data[key] = value
	return nil
}

func (i InMemoryDB) Data() map[string]model.RedirectionData {
	return i.data
}

func (i *InMemoryDB) Restore(data map[string]model.RedirectionData) {
	i.data = data
}
