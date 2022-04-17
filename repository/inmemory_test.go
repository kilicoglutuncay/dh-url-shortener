package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInMemoryRepository_Set should return nil if the key is not exists.
func TestInMemoryRepository_Set_(t *testing.T) {
	data := make(map[string]string)
	repo := InMemoryRepository{data: data}
	err := repo.Set("key", "value")
	assert.Equal(t, data["key"], "value")
	assert.Nil(t, err)
}

// TestInMemoryRepository_Set should return error if the key already exists.
func TestInMemoryRepository_Set_ShouldReturnError(t *testing.T) {
	data := make(map[string]string)
	repo := InMemoryRepository{data: data}
	_ = repo.Set("key", "value")
	err := repo.Set("key", "value2")
	assert.Equal(t, data["key"], "value")
	assert.Error(t, err)
}
