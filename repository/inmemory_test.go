package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInMemoryRepository_Set should return nil if the key is not exists.
func TestInMemoryRepository_Set_(t *testing.T) {
	data := make(map[string]string)
	repo := NewInMemoryRepository(data)
	err := repo.Set("key", "value")
	assert.Equal(t, data["key"], "value")
	assert.Nil(t, err)
}

// TestInMemoryRepository_Set should return error if the key already exists.
func TestInMemoryRepository_Set_ShouldReturnError(t *testing.T) {
	data := make(map[string]string)
	repo := NewInMemoryRepository(data)
	_ = repo.Set("key", "value")
	err := repo.Set("key", "value2")
	assert.Equal(t, data["key"], "value")
	assert.Error(t, err)
}

func TestNewInMemoryRepository(t *testing.T) {
	data := map[string]string{
		"key-1": "value-1",
		"key-2": "value-2",
	}
	r := NewInMemoryRepository(data)
	assert.Equal(t, data, r.data)
}

func TestInMemoryRepository_Get_ShouldReturnErrorWhenHashNotFound(t *testing.T) {
	data := map[string]string{
		"key-1": "value-1",
	}
	r := NewInMemoryRepository(data)
	_, err := r.Get("key-2")
	assert.Error(t, err)
}

func TestInMemoryRepository_Get(t *testing.T) {
	data := map[string]string{
		"key-1": "value-1",
	}
	r := NewInMemoryRepository(data)
	val, err := r.Get("key-1")
	assert.Nil(t, err)
	assert.Equal(t, "value-1", val)
}
