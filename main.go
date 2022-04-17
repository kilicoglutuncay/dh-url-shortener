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
	h := handler.UrlHandler{ShortenerService: shortenerService}

	s.Post("/short", h.Shorten, s.AccessLogMiddleware)

	log.Fatal(s.ListenAndServe())
}
