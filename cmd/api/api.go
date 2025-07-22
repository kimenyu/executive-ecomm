package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	apiRouter.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		err := s.db.Ping()
		if err != nil {
			http.Error(w, "DB not connected: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "API is running and DB is connected!")
	}).Methods("GET")

	log.Printf("erver listening on %s", s.addr)
	return http.ListenAndServe(s.addr, router)
}
