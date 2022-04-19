package main

import (
	"dh-url-shortener/config"
	"fmt"
	"net/http"
	"regexp"
)

var expandRe = regexp.MustCompile(`^/[a-z0-9A-Z]{7}$`)

// HTTPServer is the server that handles the HTTP requests
type HTTPServer struct {
	ServerMux  *http.ServeMux
	routeTable map[string]http.HandlerFunc
	Config     *config.Config
}

// NewHTTPServer creates a new HTTPServer
func NewHTTPServer(c *config.Config) *HTTPServer {
	mux := http.NewServeMux()
	server := &HTTPServer{
		ServerMux:  mux,
		routeTable: make(map[string]http.HandlerFunc),
		Config:     c,
	}

	return server
}

// Get is a shortcut for mapping GET requests to the specified path.
func (s *HTTPServer) Get(path string, handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	for _, m := range middlewares {
		handler = m(handler)
	}
	s.routeTable[http.MethodGet+" "+path] = handler
}

// Post is a shortcut for mapping POST requests to the specified path.
func (s *HTTPServer) Post(path string, handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) {
	for _, m := range middlewares {
		handler = m(handler)
	}
	s.routeTable[http.MethodPost+" "+path] = handler
}

// ServeHTTP routes the request to the appropriate handler
func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := s.routeTable[r.Method+" "+r.URL.Path]
	if !ok {
		if expandRe.MatchString(r.URL.Path) {
			handler = s.routeTable["GET /:hash"]
		} else {
			handler = http.NotFound
		}
	}
	handler(w, r)
}

// ListenAndServe starts the HTTP server
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
