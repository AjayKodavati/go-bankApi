package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIError struct {
  Error string
}

type APIServer struct {
	listenAddr string
	store storage
}

func NewAPIServer(listenAddr string, store storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store: store,
	}
}

func (s *APIServer) Run() {
  router := mux.NewRouter()
  
  router.HandleFunc("/account", makeHttpHandleFunc(s.handleAccount))
  router.HandleFunc("/account/{id}", makeHttpHandleFunc(s.handleGetAccountById))

  log.Println("Json Api server listening on", s.listenAddr)
  
  http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
  if r.Method == "GET" {
    return s.handleGetAccount(w, r)
  }
  if r.Method == "POST" {
    return s.handleCreateAccount(w, r)
  }
	return nil
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
  account := newAccount("Ajay", "kodavati")
  return WriteJson(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
  createAccountReq := new(CreateAccount)
  if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
    return err
  }

  account := newAccount(createAccountReq.FirstName, createAccountReq.LastName)
  if err := s.store.CreateAccount(account); err != nil {
    return err;
  }
  
  return WriteJson(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
  return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
  return nil
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
  return nil
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error 

func WriteJson(w http.ResponseWriter, statusCode int, data interface{}) error {
  w.Header().Add("Content-Type", "application/json")
  w.WriteHeader(statusCode)
  return json.NewEncoder(w).Encode(data)
}

func makeHttpHandleFunc(api apiFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    if err := api(w, r); err != nil {
      WriteJson(w, http.StatusBadRequest, APIError{Error: err.Error()})
    }
  }  
}