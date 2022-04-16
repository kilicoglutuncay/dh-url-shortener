package main

import (
	"dh-url-shortener/config"
	"dh-url-shortener/handler"
	"dh-url-shortener/service"
	"log"
)

func main() {
	c := config.NewConfig()
	s := NewHTTPServer(c)

	shortenerService := service.Shortener{}
	h := handler.ShortenerHandler{ShortenerService: shortenerService}

	s.Post("/short", h.Create, s.AccessLogMiddleware)

	log.Fatal(s.ListenAndServe())
}
