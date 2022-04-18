package main

import (
	"dh-url-shortener/config"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

type HTTPServer struct {
	ServerMux  *http.ServeMux
	routeTable map[string]http.HandlerFunc
	Config     *config.Config
}

func NewHTTPServer(c *config.Config) *HTTPServer {
	mux := http.NewServeMux()
	server := &HTTPServer{
		ServerMux:  mux,
		routeTable: make(map[string]http.HandlerFunc),
		Config:     c,
	}

	server.Get("/health", server.healthHandler, server.AccessLogMiddleware)
	return server
}

func (s *HTTPServer) Get(path string, handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	for _, m := range middlewares {
		handler = m(handler)
	}
	s.routeTable[http.MethodGet+" "+path] = handler
}

func (s *HTTPServer) Post(path string, handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	for _, m := range middlewares {
		handler = m(handler)
	}
	s.routeTable[http.MethodPost+" "+path] = handler
}

var (
	shortenRe = regexp.MustCompile(`^/short/*$`)
	expandRe  = regexp.MustCompile(`^/[a-z0-9A-Z]{7}$`)
)

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && shortenRe.MatchString(r.URL.Path):
		s.routeTable["POST /short"](w, r)
	case r.Method == http.MethodGet && expandRe.MatchString(r.URL.Path):
		s.routeTable["GET /"](w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *HTTPServer) healthHandler(w http.ResponseWriter, r *http.Request) {
	_, err := io.WriteString(w, "OK")
	if err != nil {
		log.Fatalln(err)
	}
}

func (s *HTTPServer) ListenAndServe() error {
	s.ServerMux.Handle("/", s)
	fmt.Println("listening " + s.Config.Addr)
	return http.ListenAndServe(s.Config.Addr, s.ServerMux)
}

// AccessLogMiddleware is a middleware that logs the request method, path, and the response status code.
func (s *HTTPServer) AccessLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
		s.Config.Logger.Println(r.Method + " " + r.URL.Path)
	}
}
