package db

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testSnapshotFile = "test_snapshot.db"
const testSnapshotInterval = time.Second * 5

func TestNewSnapshot_Restore_ShouldNotReturnErrorWhenSnapshotFileCantOpen(t *testing.T) {
	inMemDB := NewInMemoryDB()
	snapshot := NewSnapshot("not-existing-file-location", testSnapshotInterval)

	err := snapshot.Restore(inMemDB)
	expectedData := map[string]string{}
	assert.Nil(t, err)
	assert.Equal(t, expectedData, inMemDB.Data())
}

func TestNewSnapshot_Restore_ShouldReturnErrorWhenFileContentIsNotEncodeable(t *testing.T) {
	inMemDB := NewInMemoryDB()
	snapshot := NewSnapshot(testSnapshotFile, testSnapshotInterval)

	writeDataToSnapshot(t, []byte("not encodeable"), testSnapshotFile)
	err := snapshot.Restore(inMemDB)
	expectedData := map[string]string{}
	assert.Error(t, err)
	assert.Equal(t, expectedData, inMemDB.Data())
}

func TestSnapshot_Restore(t *testing.T) {
	inMemDB := NewInMemoryDB()
	snapshot := NewSnapshot(testSnapshotFile, testSnapshotInterval)
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}
	d, _ := json.Marshal(testData)
	writeDataToSnapshot(t, d, testSnapshotFile)
	err := snapshot.Restore(inMemDB)
	assert.Nil(t, err)
	assert.Equal(t, testData, inMemDB.Data())
}

func writeDataToSnapshot(t *testing.T, data []byte, snapshotPath string) {
	file, err := os.OpenFile(snapshotPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	assert.Nil(t, err)
	defer file.Close()

	_, err = file.Write(data)
	assert.Nil(t, err)
}
