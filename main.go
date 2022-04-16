package main

import (
	"dh-url-shortener/config"
	"dh-url-shortener/handler"
	"dh-url-shortener/service"
	"log"
	"os"
)

func main() {

	c := config.NewConfig(log.New(os.Stdout, "", log.LstdFlags))
	s := NewHTTPServer(c)

	shortenerService := service.Shortener{}
	h := handler.ShortenerHandler{ShortenerService: shortenerService}

	s.Post("/short", h.Create, s.AccessLogMiddleware)

	log.Fatal(s.ListenAndServe())
}
