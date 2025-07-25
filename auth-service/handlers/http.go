package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/lyhomyna/sf/auth-service/middleware"
	"github.com/lyhomyna/sf/auth-service/models"
	"github.com/lyhomyna/sf/auth-service/service"
	userService "github.com/lyhomyna/sf/auth-service/service/user"
)

var sessionCookieName = "session-id"

type HttpServer struct {
    http *http.Server
    Services *service.Services
}

func (s *HttpServer) Run(ctx context.Context ) error {
    mux := http.NewServeMux()

    mux.HandleFunc("/register", func(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
	    writeResponseMessage(w, http.StatusMethodNotAllowed, "Use method POST instead")
	    return
	}
	s.register(w, req)
    })
    mux.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
	    writeResponseMessage(w, http.StatusMethodNotAllowed, "Use method POST instead")
	    return
	}
	s.login(w, req)
    })
    mux.HandleFunc("/logout", func(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
	    writeResponseMessage(w, http.StatusMethodNotAllowed, "Use method GET instead")
	    return
	}
	s.logout(w, req)
    })
    mux.HandleFunc("/check-auth", func(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
	    writeResponseMessage(w, http.StatusMethodNotAllowed, "Use method GET instead")
	    return
	}
	s.checkAuth(w, req)
    })
    mux.HandleFunc("/get-user", func(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
	    writeResponseMessage(w, http.StatusMethodNotAllowed, "Use method GET instead")
	    return
	}
	s.handleGetUser(w, req)
    })
    mux.HandleFunc("/health", func(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
    })

    router := middleware.OptionsMiddleware(mux)

    s.http = &http.Server {
	Addr: ":8081",
	Handler: router,
	ReadHeaderTimeout: 5 * time.Second, // mitigate risk of Slowloris Attack
    }
    
    log.Println("HTTP server is running on port 8081")
    if err := s.http.ListenAndServe(); err != nil {
	return err
    }

    return nil 
}

func (s *HttpServer) register(w http.ResponseWriter, req *http.Request) {
    userId, errHttp := s.Services.Users.CreateUser(req)
    if errHttp != nil {
	writeResponseMessage(w, errHttp.Code, "Check logs")
	return
    }
    
    resp, err := makeCreateRootDirRequest(userId)
    if err != nil {
	log.Println(err)

	writeResponseMessage(w, http.StatusInternalServerError, "Couldn't create user")
	return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
	s.Services.Users.DeleteUser(userId)

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Root dir creation failed: %s", body)
	return
    }

    log.Printf("New user '%s' has been registered.", userId)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    response, err := json.Marshal(struct {
	Id   string `json:"id"`
    }{
	Id: userId,
    })
    if err != nil {
	panic(err)
    }
    
    w.Write(response)
}

func makeCreateRootDirRequest(userId string) (*http.Response, error) {
    serviceToken := os.Getenv("AUTH_TO_FILE_TOKEN")

    data := url.Values{}
    data.Set("userId", userId)

    requestToFilesService, err := http.NewRequest(http.MethodPost, "http://file-service:8082/create-root", strings.NewReader(data.Encode()))
    if err != nil {
	return nil, err
    }
    requestToFilesService.Header.Add("Authorization", fmt.Sprintf("Bearer %s", serviceToken))
    requestToFilesService.Header.Add("Content-Type", "application/x-www-form-urlencoded")

    client := &http.Client{
	Timeout: 5 * time.Second,
    }
    resp, err := client.Do(requestToFilesService)
    if err != nil {
	return nil, err
    }

    return resp, nil
}

func (s *HttpServer) login(w http.ResponseWriter, req *http.Request) {
    // Decode user
    var userData models.User
    err := json.NewDecoder(req.Body).Decode(&userData)
    if err != nil {
	writeResponseMessage(w, http.StatusBadRequest, "Couldn't parse user")
	return
    }

    dbUser, httpErr := s.validateUser(&userData)
    if httpErr != nil {
	writeResponseMessage(w, httpErr.Code, httpErr.Message)
	return
    }

    // create session
    sessionId, errHttp := s.Services.Sessions.CreateSession(dbUser.Id)
    if errHttp != nil {
	writeResponseMessage(w, errHttp.Code, errHttp.Message)
	return
    }
    setSessionCookie(w, sessionId)

    // write response
    log.Printf("User %s logged in.\n", dbUser.Id)
    w.WriteHeader(http.StatusOK)
}

func (s *HttpServer) validateUser(user *models.User) (*models.DbUser, *models.HTTPError) {
    if strings.Trim(user.Email, " ") == "" || strings.Trim(user.Password,  " ") == "" {
	return nil, &models.HTTPError{
	    Code: http.StatusBadRequest,
	    Message: "User data can't be blank line",
	}
    }

    if len(user.Password) < 6 {
	return nil, &models.HTTPError{
	    Code: http.StatusBadRequest,
	    Message: "Password should be at least 6 chars length",
	}
    }

    if strings.Contains(user.Password, "'\" ") {
	return nil, &models.HTTPError{
	    Code: http.StatusBadRequest,
	    Message: "Password shouldn't contain ' or \"",
	}
    }

    // find user by email
    dbUser, httpError := s.Services.Users.GetUserByEmail(user.Email);
    if httpError != nil {
	return nil, &models.HTTPError{
	    Code: http.StatusForbidden,
	    Message: "Incorrect credentials",
	}   
    }

    if err := userService.ComparePasswords(dbUser.Password, user.Password); err != nil {
	return nil, &models.HTTPError{
	    Code: http.StatusForbidden,
	    Message: "Incorrect credentials",
	}
    }

    return dbUser, nil
}

func setSessionCookie(w http.ResponseWriter, sessionId string) {
    http.SetCookie(w, &http.Cookie{
	Name: sessionCookieName,
	Value: sessionId,
	HttpOnly: true,
	Path: "/",
	SameSite: http.SameSiteLaxMode,
	MaxAge:   86400,
	Expires:  time.Now().Add(24 * time.Hour),
    })
}

func (s *HttpServer) logout(w http.ResponseWriter, req *http.Request) {
    cookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	writeResponseMessage(w, http.StatusUnauthorized, err.Error())
	return
    }
    
    sessionId := cookie.Value

    errHttp := s.Services.Sessions.DeleteSession(sessionId)
    if errHttp != nil {
	writeResponseMessage(w, errHttp.Code, errHttp.Message)
	return
    }

    http.SetCookie(w, &http.Cookie{
	Name: sessionCookieName,
	MaxAge: -1,
	Value: "",
	Path: "/",
    })	

    log.Printf("Session '%s' closed.", sessionId)
    w.WriteHeader(http.StatusOK)
}

func (s *HttpServer) checkAuth(w http.ResponseWriter, req *http.Request) {
    logConnection(req)

    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	log.Println("No cookie present")
	w.WriteHeader(http.StatusUnauthorized)
	return
    }

    if uid, ok := s.Services.Sessions.IsSessionExists(sessionCookie.Value); ok {
	w.WriteHeader(http.StatusOK)
	jsonResponse, _ := json.Marshal(map[string]string{"userId": uid})
	w.Write(jsonResponse)
	return
    }

    log.Println("WTF IS THAT COOKIE???")
    http.SetCookie(w, &http.Cookie{
	Name: sessionCookieName,
	MaxAge: -1,
	Value: "",
	Path: "/",
    })	
    w.WriteHeader(http.StatusUnauthorized)
}

func (s *HttpServer) handleGetUser(w http.ResponseWriter, req *http.Request) {
    sessionCookie, err := req.Cookie(sessionCookieName)
    if err != nil {
	writeResponseMessage(w, http.StatusUnauthorized, err.Error())
	return
    }

    userId, sessionExists := s.Services.Sessions.IsSessionExists(sessionCookie.Value)
    if !sessionExists {
	writeResponseMessage(w, http.StatusNotFound, "There is no user associated with session")
	return
    }

    user, httpErr := s.Services.Users.GetUserById(userId)
    if httpErr != nil {
	writeResponseMessage(w, httpErr.Code, httpErr.Message)
	return
    }

    // TODO: change it 
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    jsonResponse, _ := json.Marshal(struct { 
	Email string `json:"email"`
	ImageUrl string `json:"imageUrl"`
    } { Email: user.Email, ImageUrl: user.ImageUrl })
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
