package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.HandleAccount))
	router.HandleFunc("/accounts/{id}", makeHTTPHandleFunc(s.HandleGetAccountByID))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.HandleTransfer))

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
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		account, err := s.Store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	}

	if r.Method == "DELETE" {
		return s.HandleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)

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
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.Store.DeleteAcccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deletedId": id})
}

func (s *APIServer) HandleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)

	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}

	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferReq)
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
	Error string `json:"error"`
}

type APIFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, err
	}

	return id, nil
}
