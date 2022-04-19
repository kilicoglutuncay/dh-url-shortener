package db

import (
	"dh-url-shortener/internal/api/model"
	"errors"
	"sync"
)

// InMemoryDB is an in-memory implementation of the DB interface
type InMemoryDB struct {
	data  map[string]model.RedirectionData
	mutex sync.RWMutex
}

// NewInMemoryDB creates a new in-memory DB
func NewInMemoryDB() *InMemoryDB {
	repo := &InMemoryDB{
		data:  make(map[string]model.RedirectionData),
		mutex: sync.RWMutex{},
	}

	return repo
}

// Get retrieves a model.RedirectionData from the DB with the given key
func (i *InMemoryDB) Get(key string) (model.RedirectionData, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	if value, ok := i.data[key]; ok {
		return value, nil
	}
	return model.RedirectionData{}, errors.New(key + " not found")
}

// Set stores a model.RedirectionData in the DB with the given key
func (i *InMemoryDB) Set(key string, value model.RedirectionData) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if _, ok := i.data[key]; ok {
		return errors.New(key + " key already exists")
	}
	i.data[key] = value
	return nil
}

// Hit increments the hit count of the model.RedirectionData with the given key
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

// Data returns the in-memory DB data
func (i *InMemoryDB) Data() map[string]model.RedirectionData {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	return i.data
}

// Restore restores the in-memory DB data from the given data
func (i *InMemoryDB) Restore(data map[string]model.RedirectionData) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()
	i.data = data
}
