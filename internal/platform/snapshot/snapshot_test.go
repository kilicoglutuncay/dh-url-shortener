package snapshot

import (
	"dh-url-shortener/internal/api/model"
	"dh-url-shortener/internal/platform/db"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testSnapshotFile = "test_snapshot.db"
const testSnapshotInterval = time.Second * 2

func TestNewSnapshot_Restore_ShouldNotReturnErrorWhenSnapshotFileCantOpen(t *testing.T) {
	inMemDB := db.NewInMemoryDB()
	snapshot := NewSnapshot("not-existing-file-location", testSnapshotInterval)

	err := snapshot.Restore(inMemDB)
	expectedData := map[string]model.RedirectionData{}
	assert.Nil(t, err)
	assert.Equal(t, expectedData, inMemDB.Data())
}

func TestNewSnapshot_Restore_ShouldReturnErrorWhenFileContentIsNotEncodeable(t *testing.T) {
	inMemDB := db.NewInMemoryDB()
	snapshot := NewSnapshot(testSnapshotFile, testSnapshotInterval)

	writeDataToSnapshot(t, []byte("not encodeable"), testSnapshotFile)
	defer os.Truncate(testSnapshotFile, 0)
	err := snapshot.Restore(inMemDB)
	expectedData := map[string]model.RedirectionData{}
	assert.Error(t, err)
	assert.Equal(t, expectedData, inMemDB.Data())
}

func TestSnapshot_Restore(t *testing.T) {
	inMemDB := db.NewInMemoryDB()
	snapshot := NewSnapshot(testSnapshotFile, testSnapshotInterval)
	testData := map[string]model.RedirectionData{
		"key1": {OriginalURL: "value1"},
		"key2": {OriginalURL: "value2"},
		"key3": {OriginalURL: "value3"},
	}
	d, _ := json.Marshal(testData)
	writeDataToSnapshot(t, d, testSnapshotFile)
	defer os.Truncate(testSnapshotFile, 0)
	err := snapshot.Restore(inMemDB)
	assert.Nil(t, err)
	assert.Equal(t, testData, inMemDB.Data())
}

func TestSnapshot_SavePeriodically(t *testing.T) {
	inMemDB := db.NewInMemoryDB()
	inMemDB2 := db.NewInMemoryDB()
	snapshot := NewSnapshot(testSnapshotFile, testSnapshotInterval)
	testData := map[string]model.RedirectionData{
		"key1": {OriginalURL: "value1"},
		"key2": {OriginalURL: "value2"},
		"key3": {OriginalURL: "value3"},
	}
	inMemDB.Restore(testData)
	stopTimerCh := make(chan bool)
	time.AfterFunc(testSnapshotInterval+time.Second*1, func() {
		stopTimerCh <- true
	})
	snapshot.SavePeriodically(inMemDB, stopTimerCh)

	defer os.Remove(testSnapshotFile)
	inMemDB2.Restore(testData)
	assert.Equal(t, testData, inMemDB2.Data())
	assert.FileExists(t, testSnapshotFile)
}

func TestSnapshot_snapshot_ShouldReturnErrorWhenFileCanNotOpenForWrite(t *testing.T) {
	inMemDB := db.NewInMemoryDB()
	snapshot := NewSnapshot("", testSnapshotInterval)
	err := snapshot.snapshot(inMemDB)
	assert.Error(t, err)
}

func writeDataToSnapshot(t *testing.T, data []byte, snapshotPath string) {
	file, err := os.OpenFile(snapshotPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	assert.Nil(t, err)
	defer file.Close()

	_, err = file.Write(data)
	assert.Nil(t, err)
}
