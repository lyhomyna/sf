package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/lyhomyna/sf/auth-service/models"
)

func CreateUser(siglog *models.Siglog, req *http.Request) (string, *models.HTTPError) {
    defer req.Body.Close()

    var user *models.User
    if err := decodeFromTo(req.Body, user); err != nil {
	return "", &models.HTTPError {
	    Code: http.StatusBadRequest,
	    Message: "Use correct user schema.",
	} 
    }

    userId, err := siglog.Users.CreateUser(user)
    if err != nil {
	return "", &models.HTTPError {
	    Code: http.StatusInternalServerError,
	    Message: "Couldn't create user.",
	}
    }
    
    log.Println("New user and session has been created.")
    return userId, nil
}

// danger function
func DeleteUser(userId string, siglog *models.Siglog) *models.HTTPError {
    if err := siglog.Users.DeleteUser(userId); err != nil {
	return &models.HTTPError {
	    Code: http.StatusInternalServerError,
	    Message: "Couldn't delete user.",
	}
    }

    log.Println("User and session has been deleted.")
    return nil
}

func decodeFromTo(rc io.ReadCloser, target any) error {
    decoder := json.NewDecoder(rc)
    if err := decoder.Decode(target); err != nil {
        return errors.New(fmt.Sprintf("Decode failure. %s", err))
    }
    return nil
}
