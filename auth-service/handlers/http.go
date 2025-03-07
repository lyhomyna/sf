package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lyhomyna/sf/auth-service/handlers/session"
	"github.com/lyhomyna/sf/auth-service/handlers/user"
	"github.com/lyhomyna/sf/auth-service/models"
	"github.com/lyhomyna/sf/auth-service/repository"
)

var sessionCookieName = "session-id"

type httpServer struct {
    http *http.Server
}

func (s *httpServer) Run(ctx context.Context) error {
    siglog := repository.GetSiglog()
    if siglog == nil {
	return errors.New("Couldn't get SigLog.")
    }

    mux := http.NewServeMux()

    mux.HandleFunc("/register", func(w http.ResponseWriter, req *http.Request) {
	register(siglog, w, req)
    })
    mux.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
	login(siglog, w, req)
    })
    mux.HandleFunc("/logout", func(w http.ResponseWriter, req *http.Request) {
	logout(siglog, w, req)
    })

    s.http = &http.Server {
	Addr: ":8080",
	Handler: mux,
	ReadHeaderTimeout: 5 * time.Second, // mitigate risk of Slowloris Attack
    }
    
    log.Println("HTTP server is running on port 8080")
    if err := s.http.ListenAndServe(); err != nil {
	return err
    }

    return nil 
}

func register(siglog *models.Siglog, w http.ResponseWriter, req *http.Request) {
    userId, errHttp := user.CreateUser(siglog, req)
    if errHttp != nil {
	return
    }

    sessionId, errHttp := session.Create(userId, siglog)
    if errHttp != nil {
	writeResponse(w, errHttp.Code, errHttp.Message)
	return
    }

    http.SetCookie(w, &http.Cookie{
	Name: sessionCookieName,
	Value: sessionId,
    })	

    log.Println("Session created")
    w.WriteHeader(http.StatusOK)
}

func login(siglog *models.Siglog, w http.ResponseWriter, req *http.Request) {
    // validate user (validate by input data and if user exists in DB)
    var inputData struct {
	Email string 	`json:"email"`
	Password string	`json:"pwd"`
    }
    err := json.NewDecoder(req.Body).Decode(inputData)
    if err != nil {
	writeResponse(w, http.StatusBadRequest, "Couldn't parse user.")
	return
    }

    if strings.Trim(inputData.Email, " ") == "" || strings.Trim(inputData.Password,  " ") == "" {
	writeResponse(w, http.StatusBadRequest, "User data can't be blank line.")
	return
    }

    if len(inputData.Password) < 6 {
	writeResponse(w, http.StatusBadRequest, "Password should be at least 6 char length.")
	return
    }

    if strings.Contains(inputData.Password, "'\" ") {
	writeResponse(w, http.StatusBadRequest, "Password shouldn't contain ' or \".")
	return
    }

    

    // create session
    // write response
    panic("Not yet implemented.")
}

func logout(siglog *models.Siglog, w http.ResponseWriter, req *http.Request) {
    cookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	writeResponse(w, http.StatusUnauthorized, err.Error())
	return
    }
    
    err = siglog.Sessions.DeleteSession(cookie.Value)
    if err != nil {
	writeResponse(w, http.StatusInternalServerError, err.Error())
	return
    }

    http.SetCookie(w, &http.Cookie{
	Name: sessionCookieName,
	MaxAge: -1,
    })	

    w.WriteHeader(http.StatusOK)
}

func writeResponse(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    responseMessage := fmt.Sprintf("'message':'%s'", message)
    w.Write([]byte(responseMessage))
}
