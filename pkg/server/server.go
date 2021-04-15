package server

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/CalvoM/file-upload/pkg/auth"
	"github.com/CalvoM/file-upload/pkg/models"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	pg_models "github.com/go-oauth2/oauth2/v4/models"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type ClientCred struct {
	ID     string `json:"client_id"`
	Secret string `json:"client_secret"`
}

func GetNewServer() (srv *http.Server) {
	auth.TokenServer.SetPasswordAuthorizationHandler(AuthPasswordHandler)
	auth.TokenServer.SetInternalErrorHandler(AuthInternalErrorHandler)
	auth.TokenServer.SetResponseErrorHandler(AuthResponseErrorHandler)
	r := mux.NewRouter()
	r.Path("/oauth/client_cred/").HandlerFunc(clientCredentialsHandler).Methods("GET")
	r.Path("/oauth/token/").HandlerFunc(tokenHandler).Methods("POST")
	r.Path("/signup/").HandlerFunc(isAuthorized(signUpHandler)).Methods("POST")
	r.Path("/upload_file/").HandlerFunc(isAuthorized(fileHandler)).Methods("POST")
	return &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:4009",
	}
}

func clientCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	res := ClientCred{}
	res.ID, res.Secret = getCredentials()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	err := auth.TokenServer.HandleTokenRequest(w, r)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"username"`
	}
	res := resp{}
	err := json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user := models.User{
		Email:    res.Email,
		Password: res.Password,
		UserName: res.Name,
	}
	_, err = user.AddUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%s", "Added")

}

func AuthPasswordHandler(username, password string) (userID string, err error) {
	user := models.User{}
	err = user.GetUser(username, password)
	if err != nil {
		log.Error(err)
		return
	}
	userID = strconv.Itoa(int(uint64(user.ID)))
	return
}

func AuthInternalErrorHandler(err error) (re *errors.Response) {
	log.Error(err)
	return

}

func AuthResponseErrorHandler(re *errors.Response) {
	log.Error(re.Error.Error())
}

func getCredentials() (string, string) {
	clientID := generateSecurityCredentials(18)
	clientSecret := generateSecurityCredentials(36)
	var client oauth2.ClientInfo = &pg_models.Client{
		ID:     clientID,
		Secret: clientSecret,
		Domain: "http://localhost:5000",
		UserID: clientID,
	}
	err := auth.ClientStore.Create(client)
	if err != nil {
		log.Fatal(err.Error())
	}
	auth.TokenManager.MapClientStorage(auth.ClientStore)
	return clientID, clientSecret
}
func generateSecurityCredentials(size int) string {
	key := make([]byte, size)
	_, _ = rand.Read(key[:])
	return base64.URLEncoding.EncodeToString(key[:])
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.TokenServer.ValidationBearerToken(r)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		endpoint(w, r)
	})
}
