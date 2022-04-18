package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInMemoryRepository_Set should return nil if the key is not exists.
func TestInMemoryRepository_Set(t *testing.T) {
	inMemoryDB := NewInMemoryDB()
	err := inMemoryDB.Set("key", "value")
	assert.Equal(t, inMemoryDB.data["key"], "value")
	assert.Nil(t, err)
}

// TestInMemoryRepository_Set should return error if the key already exists.
func TestInMemoryRepository_Set_ShouldReturnErrorWhenKeyAlreadyExists(t *testing.T) {
	inMemoryDB := NewInMemoryDB()
	_ = inMemoryDB.Set("key", "value1")
	err := inMemoryDB.Set("key", "value2")
	assert.Equal(t, inMemoryDB.data["key"], "value1")
	assert.Error(t, err)
}

func TestNewInMemoryDB_ShouldReturnRepoWithSavedSnapshotDb(t *testing.T) {
	testData := map[string]string{"05bf184": "https://www.yemeksepeti.com/istanbul"}
	inMemoryDB := NewInMemoryDB()
	inMemoryDB.Restore(testData)
	assert.Equal(t, inMemoryDB.data, testData)
}

func TestInMemoryRepository_Get_ShouldReturnErrorWhenHashNotFound(t *testing.T) {
	inMemoryDB := NewInMemoryDB()
	inMemoryDB.data["key-1"] = "value-1"
	_, err := inMemoryDB.Get("key-2")
	assert.Error(t, err)
}

func TestInMemoryRepository_Get(t *testing.T) {
	inMemoryDB := NewInMemoryDB()
	inMemoryDB.data["key-1"] = "value-1"
	val, err := inMemoryDB.Get("key-1")
	assert.Nil(t, err)
	assert.Equal(t, "value-1", val)
}

func TestInMemoryDB_Data(t *testing.T) {
	inMemoryDB := NewInMemoryDB()
	inMemoryDB.data["key-1"] = "value-1"
	assert.Equal(t, inMemoryDB.data, inMemoryDB.Data())
}
