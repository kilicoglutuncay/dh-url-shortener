package main

import (
	"dh-url-shortener/config"
	"dh-url-shortener/internal/api/handler"
	"dh-url-shortener/internal/api/service"
	"dh-url-shortener/internal/platform/db"
	dbSnapshot "dh-url-shortener/internal/platform/snapshot"
	"fmt"
	"log"
	"os"
)

func main() {
	c := config.NewConfig(log.New(os.Stdout, "", log.LstdFlags))
	fmt.Printf("Config: %#v\n", c)
	s := NewHTTPServer(c)
	inMemoryDB := db.NewInMemoryDB()
	snapshot := dbSnapshot.NewSnapshot(c.DBSnapshotPath, c.SnapshotSaveInterval)
	err := snapshot.Restore(inMemoryDB)
	if err != nil {
		log.Fatal(err)
	}
	go snapshot.SavePeriodically(inMemoryDB)

	shortenerService := service.Shortener{DB: inMemoryDB, ShortURLDomain: c.ShortURLDomain}
	h := handler.URLHandler{ShortenerService: shortenerService}

	s.Post("/short", h.Shorten, s.AccessLogMiddleware)
	s.Get("/:hash", h.Expand, s.AccessLogMiddleware)
	s.Get("/list", h.List, s.AccessLogMiddleware)

	log.Fatal(s.ListenAndServe())
}
