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
	    Message: err.Error(),
	}
    }
    log.Println("Session created")
    return sessionId, nil
}

func Delete(sessionId string, siglog *models.Siglog) *models.HTTPError {
    err := siglog.Sessions.DeleteSession(sessionId)
    if err != nil {
	return &models.HTTPError{
	    Code: http.StatusInternalServerError,
	    Message: err.Error(),
	}
    }
    return nil
}

func IsSessionExists(sessionId string, siglog *models.Siglog) bool {
    userId, err := siglog.Sessions.UserIdFromSessionId(sessionId)
    log.Printf("UserID for session %s is %s\n", sessionId, userId)
    log.Printf("FUCKING ERROR is: %s\n", err.Error())
    if userId != "" {
	return true 
    }
    return false
}
