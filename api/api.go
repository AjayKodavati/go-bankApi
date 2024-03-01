package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type APIError struct {
	Error string `json:"error"`
}

type APIServer struct {
	listenAddr string
	store      storage
}

func NewAPIServer(listenAddr string, store storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHttpHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHttpHandleFunc(s.handleGetAccountById), s.store))
	router.HandleFunc("/transfer", makeHttpHandleFunc(s.handleTransfer))
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
		return err
	}

	tokenString, err := createJwt(account)
	if err != nil {
		return err
	}
	fmt.Print("token string ", tokenString)

	return WriteJson(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err = s.store.DeleteAccount(id); err != nil {
		return fmt.Errorf("unable to delete account %d", id)
	}
	return WriteJson(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transfer := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transfer); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJson(w, http.StatusOK, transfer)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}
		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}
		return WriteJson(w, http.StatusOK, account)
	} else if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("unsupported method")
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

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}

	return id, nil
}

func permissionDenied(w http.ResponseWriter) {
	 WriteJson(w, http.StatusForbidden, APIError{Error: "permisson denied"})
}

func withJWTAuth(handlerFunc http.HandlerFunc, s storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokeString := r.Header.Get(("x-jwt-token"))

		token, err := validateToken(tokeString)
		fmt.Printf("token -> %s\n", token.Claims.Valid())
		if err != nil {
			permissionDenied(w)
			return
		}
		if !token.Valid {
			permissionDenied(w)
			return
		}

		userId, err := getID(r)

		if err != nil {
			permissionDenied(w)
			return 
		}

		account, err := s.GetAccountByID(userId)
		if err != nil {
			permissionDenied(w)
			return
		}
		if account.Number != int64(token.Claims.(jwt.MapClaims)["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}
		handlerFunc(w, r)
	}
}

func validateToken(tokenString string) (*jwt.Token, error) {
	secret := "hunter99"
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

func createJwt(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"ExpiresAt": 15000,
		"AccountNumber": account.Number,
	}

	secret := "hunter99"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}