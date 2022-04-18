package main

import (
	"bytes"
	"dh-url-shortener/config"
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

	r, _ := http.NewRequest("GET", "http://localhost:8080/shorten", nil)
	w := httptest.NewRecorder()
	handler(w, r)
	line := buf.String()
	expectedLog := "GET /shorten"

	assert.True(t, strings.Contains(line, expectedLog))
}
