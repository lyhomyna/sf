package handlers 

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/lyhomyna/sf/auth-service/service/session"
	"github.com/lyhomyna/sf/auth-service/service/user"
	"github.com/lyhomyna/sf/auth-service/models"
	"github.com/lyhomyna/sf/auth-service/repository"
)

var sessionCookieName = "session-id"

type HttpServer struct {
    http *http.Server
}

func (s *HttpServer) Run(ctx context.Context) error {
    siglog := repository.GetSiglog()
    if siglog == nil {
	return errors.New("Couldn't get SigLog.")
    }

    mux := http.NewServeMux()

    mux.HandleFunc("/register", func(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
	    writeResponseMessage(w, http.StatusMethodNotAllowed, "Use method POST instead")
	    return
	}
	register(siglog, w, req)
    })
    mux.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
	    writeResponseMessage(w, http.StatusMethodNotAllowed, "Use method POST instead")
	    return
	}
	login(siglog, w, req)
    })
    mux.HandleFunc("/logout", func(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
	    writeResponseMessage(w, http.StatusMethodNotAllowed, "Use method GET instead")
	    return
	}
	logout(siglog, w, req)
    })
    mux.HandleFunc("/check-auth", func(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
	    writeResponseMessage(w, http.StatusMethodNotAllowed, "Use method GET instead")
	    return
	}
	checkAuth(siglog, w, req)
    })
    mux.HandleFunc("/get-user", func(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
	    writeResponseMessage(w, http.StatusMethodNotAllowed, "Use method GET instead")
	    return
	}
	handleGetUser(siglog, w, req)
    })
    mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
    })

    s.http = &http.Server {
	Addr: ":8081",
	Handler: mux,
	ReadHeaderTimeout: 5 * time.Second, // mitigate risk of Slowloris Attack
    }
    
    log.Println("HTTP server is running on port 8081")
    if err := s.http.ListenAndServe(); err != nil {
	return err
    }

    return nil 
}

func register(siglog *models.Siglog, w http.ResponseWriter, req *http.Request) {
    userId, errHttp := user.CreateUser(siglog, req)
    if errHttp != nil {
	writeResponseMessage(w, errHttp.Code, errHttp.Message)
	return
    }
    
    sessionId, errHttp := session.Create(userId, siglog)
    if errHttp != nil {
	writeResponseMessage(w, errHttp.Code, errHttp.Message)
	return
    }

    http.SetCookie(w, &http.Cookie{
	Name: sessionCookieName,
	Value: sessionId,
    })	

    log.Printf("New user '%s' has been registered.", userId)
    w.WriteHeader(http.StatusOK)
}

func login(siglog *models.Siglog, w http.ResponseWriter, req *http.Request) {
    // validate user (validate by input data and if user exists in DB)
    var userData models.User
    err := json.NewDecoder(req.Body).Decode(&userData)
    if err != nil {
	writeResponseMessage(w, http.StatusBadRequest, "Couldn't parse user")
	return
    }

    if strings.Trim(userData.Email, " ") == "" || strings.Trim(userData.Password,  " ") == "" {
	writeResponseMessage(w, http.StatusBadRequest, "User data can't be blank line")
	return
    }

    if len(userData.Password) < 6 {
	writeResponseMessage(w, http.StatusBadRequest, "Password should be at least 6 chars length")
	return
    }

    if strings.Contains(userData.Password, "'\" ") {
	writeResponseMessage(w, http.StatusBadRequest, "Password shouldn't contain ' or \"")
	return
    }

    // find user by email
    dbUser, httpError := user.GetUserByEmail(userData.Email, siglog);
    if httpError != nil {
	writeResponseMessage(w, httpError.Code, httpError.Message);
	return
    }

    if err := user.ComparePasswords(dbUser.Password, userData.Password); err != nil {
	writeResponseMessage(w, http.StatusForbidden, "Passwords don't match")
	return
    }

    // create session
    sessionId, errHttp := session.Create(dbUser.Id, siglog)
    if errHttp != nil {
	writeResponseMessage(w, errHttp.Code, errHttp.Message)
	return
    }

    http.SetCookie(w, &http.Cookie{
	Name: sessionCookieName,
	Value: sessionId,
	HttpOnly: true,
	Path: "/",
	SameSite: http.SameSiteLaxMode,
    })

    // write response
    log.Printf("User %s logged in.\n", dbUser.Id)
    w.WriteHeader(http.StatusOK)
}

func logout(siglog *models.Siglog, w http.ResponseWriter, req *http.Request) {
    cookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	writeResponseMessage(w, http.StatusUnauthorized, err.Error())
	return
    }
    
    sessionId := cookie.Value

    errHttp := session.Delete(sessionId, siglog)
    if errHttp != nil {
	writeResponseMessage(w, errHttp.Code, errHttp.Message)
	return
    }

    http.SetCookie(w, &http.Cookie{
	Name: sessionCookieName,
	MaxAge: -1,
    })	

    log.Printf("Session '%s' closed.", sessionId)
    w.WriteHeader(http.StatusOK)
}

func checkAuth(siglog *models.Siglog, w http.ResponseWriter, req *http.Request) {
    logConnection(req)

    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	log.Println("No cookie present")
	w.WriteHeader(http.StatusUnauthorized)
	return
    }

    if uid, ok := session.IsSessionExists(sessionCookie.Value, siglog); ok {
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(map[string]string{"userId": uid})
	w.Write(jsonResponse)
	return
    }

    log.Println("WTF IS THAT COOKIE???")
    http.SetCookie(w, &http.Cookie{
	Name: sessionCookieName,
	MaxAge: -1,
    })	
    w.WriteHeader(http.StatusUnauthorized)
}

func handleGetUser(siglog *models.Siglog, w http.ResponseWriter, req *http.Request) {
    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	writeResponseMessage(w, http.StatusUnauthorized, err.Error())
	return
    }

    userId, sessionExists := session.IsSessionExists(sessionCookie.Value, siglog)
    if !sessionExists {
	writeResponseMessage(w, http.StatusNotFound, "There is no associated user with session")
	return
    }

    user, httpErr := user.GetById(userId, siglog)
    if err != nil {
	writeResponseMessage(w, httpErr.Code, httpErr.Message)
	return
    }

    // TODO: change it 
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    jsonResponse, _ := json.Marshal(struct { 
	Email string `json:"email"`
    } { Email: user.Email })
    w.Write(jsonResponse)
}

func writeResponseMessage(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    jsonResponse, _ := json.Marshal(map[string]string{"message": message})
    w.Write(jsonResponse)
}

// logConnection is a 'logger' of each request
func logConnection(req *http.Request) {
    log.Printf("%s | %s %s\n", req.RemoteAddr, req.Method,  req.URL.Path)
}
