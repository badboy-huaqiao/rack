package pkg

import (
	"net/http"
	"rack/pkg/handler"

	"github.com/gorilla/mux"
)

func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/plugin/register", handler.Register).Methods(http.MethodPost)
	s.HandleFunc("/plugin/upload", handler.Upload).Methods(http.MethodPost)
	s.HandleFunc("/plugin/load", handler.Load).Methods(http.MethodPost)
	return r
}
