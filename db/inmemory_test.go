package db

import (
	"dh-url-shortener/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInMemoryRepository_Set should return nil if the key is not exists.
func TestInMemoryRepository_Set(t *testing.T) {
	inMemoryDB := NewInMemoryDB()
	err := inMemoryDB.Set("key", model.RedirectionData{OriginalURL: "value"})
	assert.Equal(t, inMemoryDB.data["key"].OriginalURL, "value")
	assert.Nil(t, err)
}

// TestInMemoryRepository_Set should return error if the key already exists.
func TestInMemoryRepository_Set_ShouldReturnErrorWhenKeyAlreadyExists(t *testing.T) {
	inMemoryDB := NewInMemoryDB()
	_ = inMemoryDB.Set("key", model.RedirectionData{OriginalURL: "value1"})
	err := inMemoryDB.Set("key", model.RedirectionData{OriginalURL: "value2"})
	assert.Equal(t, inMemoryDB.data["key"].OriginalURL, "value1")
	assert.Error(t, err)
}

func TestNewInMemoryDB_ShouldReturnRepoWithSavedSnapshotDb(t *testing.T) {
	testData := map[string]model.RedirectionData{"05bf184": {OriginalURL: "https://www.yemeksepeti.com/istanbul"}}
	inMemoryDB := NewInMemoryDB()
	inMemoryDB.Restore(testData)
	assert.Equal(t, inMemoryDB.data, testData)
}

func TestInMemoryRepository_Get_ShouldReturnErrorWhenHashNotFound(t *testing.T) {
	inMemoryDB := NewInMemoryDB()
	inMemoryDB.data["key-1"] = model.RedirectionData{OriginalURL: "value-1"}
	_, err := inMemoryDB.Get("key-2")
	assert.Error(t, err)
}

func TestInMemoryRepository_Get(t *testing.T) {
	inMemoryDB := NewInMemoryDB()
	inMemoryDB.data["key-1"] = model.RedirectionData{OriginalURL: "value-1"}
	val, err := inMemoryDB.Get("key-1")
	assert.Nil(t, err)
	assert.Equal(t, "value-1", val.OriginalURL)
}

func TestInMemoryDB_Data(t *testing.T) {
	inMemoryDB := NewInMemoryDB()
	inMemoryDB.data["key-1"] = model.RedirectionData{OriginalURL: "value-1"}
	assert.Equal(t, inMemoryDB.data, inMemoryDB.Data())
}

// TestInMemoryRepository_Hit should return error if the key not exists.
func TestInMemoryRepository_Hit_ShouldReturnErrorWhenKeyNotExists(t *testing.T) {
	inMemoryDB := NewInMemoryDB()
	err := inMemoryDB.Hit("key")

	assert.Error(t, err)
}

// TestInMemoryRepository_Hit should return error if the key not exists.
func TestInMemoryRepository_Hit_ShouldIncreaseHitOfRedirectionData(t *testing.T) {
	inMemoryDB := NewInMemoryDB()
	inMemoryDB.data["key"] = model.RedirectionData{OriginalURL: "value1", Hits: 0}
	err := inMemoryDB.Hit("key")

	assert.Nil(t, err)
	assert.Equal(t, 1, inMemoryDB.data["key"].Hits)
}
