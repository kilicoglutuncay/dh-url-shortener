package snapshot

import (
	"dh-url-shortener/internal/api/model"
	"dh-url-shortener/internal/api/service"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Snapshot struct {
	SnapshotPath         string
	SnapshotSaveInterval time.Duration
}

func NewSnapshot(snapshotPath string, snapshotSaveInterval time.Duration) *Snapshot {
	return &Snapshot{
		SnapshotPath:         snapshotPath,
		SnapshotSaveInterval: snapshotSaveInterval,
	}
}

func (s Snapshot) snapshot(db service.DB) error {
	file, err := os.OpenFile(s.SnapshotPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(db.Data())
	if err != nil {
		return err
	}

	return nil
}

func (s Snapshot) Restore(db service.DB) error {
	file, err := os.OpenFile(s.SnapshotPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("Snapshot file not found, starting from empty database")
		return nil
	}
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	var data map[string]model.RedirectionData
	if len(byteValue) > 0 {
		if err := json.Unmarshal(byteValue, &data); err != nil {
			return err
		}
		db.Restore(data)
	}

	return nil
}

func (s Snapshot) SavePeriodically(db service.DB) {
	ticker := time.NewTicker(s.SnapshotSaveInterval)

	for { //nolint:gosimple
		select {
		case <-ticker.C:
			err := s.snapshot(db)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}
}
