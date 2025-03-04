package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/lyhomyna/sf/auth-service/database/models"
)

func Routes(siglog models.Siglog) http.Handler {
    mux := http.NewServeMux()

    mux.HandleFunc("/create", func(w http.ResponseWriter, req *http.Request) {
	createUser(siglog, w, req)
    })
    mux.HandleFunc("/delete", func(w http.ResponseWriter, req *http.Request) {
	deleteUser(siglog, w, req)
    })

    return mux
}

// POST
func createUser(siglog models.Siglog, w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if req.Method != http.MethodPost {
	writeResponse(w, http.StatusMethodNotAllowed, "Use method POST instead.")
	return
    }

    defer req.Body.Close()

    var user *models.User
    if err := decodeFromTo(req.Body, user); err != nil {
	writeResponse(w, http.StatusBadRequest, "Use correct user schema.")
	return
    }

    userId, err := siglog.Users.CreateUser(user)
    if err != nil {
	writeResponse(w, http.StatusInternalServerError, "Couldn't create user.")
	return
    }

    sessionId, err := siglog.Sessions.CreateSession(userId)
    if err != nil {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("'message': 'User created, but couldn't create session in DB.'"))
	return
    }

    http.SetCookie(w, &http.Cookie {
        Name: "session-id",
        Value: sessionId,
    })
    
    log.Println("New user and session has been created.")
    w.WriteHeader(http.StatusOK)
}

// Danger method 
// DELETE
func deleteUser(siglog models.Siglog, w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    if req.Method != http.MethodDelete {
	writeResponse(w, http.StatusMethodNotAllowed, "Use method DELETE instead.")
	return
    }

    userId := strings.TrimPrefix(req.URL.Path, "/delete/")
    if err := siglog.Users.DeleteUser(userId); err != nil {
	writeResponse(w, http.StatusInternalServerError, "Coulnd't delete user.")
	return
    }

    sessionCookie, err := req.Cookie("session-id")
    if err != nil {
	writeResponse(w, http.StatusInternalServerError, "Couldn't get cookie.")
	return
    }

    err = siglog.Sessions.DeleteSession(sessionCookie.Value)
    if err != nil {
	writeResponse(w, http.StatusInternalServerError, "Coookie not deleted.")
	return
    }

    http.SetCookie(w, &http.Cookie {
	Name: "session-id",
	MaxAge: -1,
    })

    log.Println("User and session has been deleted.")
    w.WriteHeader(http.StatusOK)
}

func decodeFromTo(rc io.ReadCloser, target any) error {
    decoder := json.NewDecoder(rc)
    if err := decoder.Decode(target); err != nil {
        return errors.New(fmt.Sprintf("Decode failure. %s", err))
    }
    return nil
}

func writeResponse(w http.ResponseWriter, code int, message string) {
    w.WriteHeader(code)
    responseMessage := fmt.Sprintf("'message':'%s'", message)
    w.Write([]byte(responseMessage))
}
