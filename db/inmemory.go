package db

import (
	"dh-url-shortener/model"
	"errors"
	"sync"
)

type InMemoryDB struct {
	data  map[string]model.RedirectionData
	mutex sync.RWMutex
}

func NewInMemoryDB() *InMemoryDB {
	repo := &InMemoryDB{
		data:  make(map[string]model.RedirectionData),
		mutex: sync.RWMutex{},
	}

	return repo
}

func (i *InMemoryDB) Get(key string) (model.RedirectionData, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	if value, ok := i.data[key]; ok {
		return value, nil
	}
	return model.RedirectionData{}, errors.New(key + " not found")
}

func (i *InMemoryDB) Set(key string, value model.RedirectionData) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if _, ok := i.data[key]; ok {
		return errors.New(key + " key already exists")
	}
	i.data[key] = value
	return nil
}

func (i *InMemoryDB) Hit(key string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	value, ok := i.data[key]
	if !ok {
		return errors.New(key + " key not exists")
	}
	value.Hits++
	i.data[key] = value
	return nil
}

func (i *InMemoryDB) Data() map[string]model.RedirectionData {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.data
}

func (i *InMemoryDB) Restore(data map[string]model.RedirectionData) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	i.data = data
}
