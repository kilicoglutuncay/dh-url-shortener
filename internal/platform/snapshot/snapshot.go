package snapshot

import (
	"dh-url-shortener/internal/api/model"
	"dh-url-shortener/internal/api/service"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

// Snapshot saves and restores the state of the database.
type Snapshot struct {
	SnapshotPath         string
	SnapshotSaveInterval time.Duration
}

// NewSnapshot creates a new snapshot object.
func NewSnapshot(snapshotPath string, snapshotSaveInterval time.Duration) *Snapshot {
	return &Snapshot{
		SnapshotPath:         snapshotPath,
		SnapshotSaveInterval: snapshotSaveInterval,
	}
}

// Save saves the current state of the database to SnapshotPath.
func (s Snapshot) snapshot(db service.DB) error {
	file, err := os.OpenFile(s.SnapshotPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_ = json.NewEncoder(file).Encode(db.Data())

	return nil
}

// Restore restores the state of the database from SnapshotPath.
func (s Snapshot) Restore(db service.DB) error {
	file, err := os.OpenFile(s.SnapshotPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("Snapshot file not found, starting from empty database")
		return nil
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	var data map[string]model.RedirectionData
	if len(byteValue) > 0 {
		if err := json.Unmarshal(byteValue, &data); err != nil {
			return err
		}
		db.Restore(data)
	}

	return nil
}

// SavePeriodically saves the state of the database within each SnapshotSaveInterval.
func (s Snapshot) SavePeriodically(db service.DB, stop chan bool) {
	ticker := time.NewTicker(s.SnapshotSaveInterval)

	for {
		select {
		case <-ticker.C:
			err := s.snapshot(db)
			if err != nil {
				log.Fatalln(err)
			}
		case <-stop:
			ticker.Stop()
			return
		}
	}
}
