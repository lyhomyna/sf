package session

import (
	"log"
	"net/http"

	"github.com/lyhomyna/sf/auth-service/models"
)

func Create(userId string, siglog *models.Siglog) (string, *models.HTTPError) {
    sessionId, err := siglog.Sessions.CreateSession(userId)
    if err != nil {
	return "", &models.HTTPError{
	    Code: http.StatusInternalServerError,
	    Message: "User created, but couldn't create session.",
	}
    }

    log.Println("Session created")
    return sessionId, nil
}

func Delete(sessionId string) {
    panic("Not yet implemented.")
}
