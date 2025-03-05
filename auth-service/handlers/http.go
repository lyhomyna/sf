package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/lyhomyna/sf/auth-service/handlers/session"
	"github.com/lyhomyna/sf/auth-service/handlers/user"
	"github.com/lyhomyna/sf/auth-service/models"
	"github.com/lyhomyna/sf/auth-service/repository"
)

var cookieSeessionIdName = "session-id"

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
	Name: cookieSeessionIdName,
	Value: sessionId,
    })	

    log.Println("Session created")
    w.WriteHeader(http.StatusOK)
}

func login(siglog *models.Siglog, w http.ResponseWriter, req *http.Request) {
    panic("Not yet implemented.")
}

func logout(siglog *models.Siglog, w http.ResponseWriter, req *http.Request) {
    panic("Not yet implemented.")
}

func writeResponse(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    responseMessage := fmt.Sprintf("'message':'%s'", message)
    w.Write([]byte(responseMessage))
}
