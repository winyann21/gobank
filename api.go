package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	// TODO: add DB
	ListenAddr string
	Store      Storage
}

func NewAPIServer(listedAddr string, store Storage) *APIServer {
	return &APIServer{
		ListenAddr: listedAddr,
		Store:      store,
	}
}

// ROUTER
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.HandleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.HandleGetAccountByID))

	log.Println("JSON API server running on port: ", s.ListenAddr)

	http.ListenAndServe(s.ListenAddr, router)
}

// ACCOUNT FUNCS
func (s *APIServer) HandleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.HandleGetAccounts(w, r)
	case "POST":
		return s.HandleCreateAccount(w, r)
	case "DELETE":
		return s.HandleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) HandleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.Store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) HandleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Println(id)

	return WriteJSON(w, http.StatusOK, &Account{})
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accReq := new(CreateAccountRequest)

	if err := json.NewDecoder(r.Body).Decode(accReq); err != nil {
		return err
	}

	account := NewAccount(accReq.FirstName, accReq.LastName)
	if err := s.Store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) HandleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

//**************Helper funcs**************

// Return as JSON
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
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
