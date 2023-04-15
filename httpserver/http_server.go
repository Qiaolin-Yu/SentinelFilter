package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sentinel-filter/service"
)

type Server struct {
	Filter *service.BloomFilter
}

func NewServer(filter *service.BloomFilter) *Server {
	return &Server{Filter: filter}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleGet(w, r)
	case http.MethodPost:
		s.handlePost(w, r)
	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Key is missing", http.StatusBadRequest)
		return
	}
	ok := s.Filter.Check(key)
	fmt.Fprint(w, ok)
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Key string `json:"key"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if data.Key == "" {
		http.Error(w, "Key is missing", http.StatusBadRequest)
		return
	}
	s.Filter.Add(data.Key)
	w.WriteHeader(http.StatusNoContent)
}
