package main

import (
	"dh-url-shortener/config"
	"dh-url-shortener/db"
	"dh-url-shortener/handler"
	"dh-url-shortener/service"
	"log"
	"os"
)

func main() {
	c := config.NewConfig(log.New(os.Stdout, "", log.LstdFlags))
	s := NewHTTPServer(c)
	inMemoryDB := db.NewInMemoryDB()
	snapshot := db.NewSnapshot(c.DBSnapshotPath, c.SnapshotSaveInterval)
	err := snapshot.Restore(inMemoryDB)
	if err != nil {
		log.Fatal(err)
	}
	go snapshot.SavePeriodically(inMemoryDB)

	shortenerService := service.Shortener{DB: inMemoryDB, ShortURLDomain: c.ShortURLDomain}
	h := handler.URLHandler{ShortenerService: shortenerService}

	s.Post("/short", h.Shorten, s.AccessLogMiddleware)
	s.Get("/", h.Expand, s.AccessLogMiddleware)
	s.Get("/list", h.List, s.AccessLogMiddleware)

	log.Fatal(s.ListenAndServe())
}
