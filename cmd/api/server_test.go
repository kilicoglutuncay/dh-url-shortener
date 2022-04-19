package main

import (
	"bytes"
	"dh-url-shortener/config"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpServer_AccessLogMiddleware(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	c := config.NewConfig(logger)
	s := NewHTTPServer(c)

	nextHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	handler := s.AccessLogMiddleware(nextHandler)

	r, _ := http.NewRequest("GET", "http://localhost:8080/shorten", http.NoBody)
	w := httptest.NewRecorder()
	handler(w, r)
	line := buf.String()
	expectedLog := "GET /shorten"

	assert.True(t, strings.Contains(line, expectedLog))
}

func TestHTTPServer_Get(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	c := config.NewConfig(logger)
	s := NewHTTPServer(c)

	handler := func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "Hello World") }
	s.Get("/test", handler, s.AccessLogMiddleware)
	_, ok := s.routeTable["GET /test"]

	assert.True(t, ok)
}

func TestHTTPServer_Post(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	c := config.NewConfig(logger)
	s := NewHTTPServer(c)

	handler := func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "Hello World") }
	s.Post("/test", handler, s.AccessLogMiddleware)
	_, ok := s.routeTable["POST /test"]

	assert.True(t, ok)
}

func TestHTTPServer_ServeHTTP_ShouldReturnStatusNotFoundWhenHandlerNotFoundForGivenEndpoint(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	c := config.NewConfig(logger)
	s := NewHTTPServer(c)

	r, _ := http.NewRequest("GET", "http://localhost:8080/not-existing-endpoint", http.NoBody)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHTTPServer_ServeHTTP_ShouldHandleDynamicHashVariable(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	c := config.NewConfig(logger)
	s := NewHTTPServer(c)
	s.Get("/:hash", func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "Hello World") })
	r, _ := http.NewRequest("GET", "http://localhost:8080/sevenCh", http.NoBody)
	w := httptest.NewRecorder()
	s.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Hello World", w.Body.String())
}
