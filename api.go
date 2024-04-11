package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Return as JSON
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// helper for handling HTTP requests
type APIError struct {
	Error string
}

type APIFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	// TODO: add DB
	ListenAddr string
}

func NewAPIServer(listedAddr string) *APIServer {
	return &APIServer{
		ListenAddr: listedAddr,
	}
}

// ROUTER
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.HandleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.HandleGetAccount))

	log.Println("JSON API server running on port: ", s.ListenAddr)

	http.ListenAndServe(s.ListenAddr, router)
}

// ACCOUNT FUNCS
func (s *APIServer) HandleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.HandleGetAccount(w, r)
	case "POST":
		return s.HandleCreateAccount(w, r)
	case "DELETE":
		return s.HandleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) HandleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Println(id)
	// account := NewAccount("Sherwin", "Romero")
	return WriteJSON(w, http.StatusOK, &Account{})
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) HandleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
