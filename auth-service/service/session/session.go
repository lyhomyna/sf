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

func IsSessionExists(sessionId string, siglog *models.Siglog) (string, bool) {
    userId, _ := siglog.Sessions.UserIdFromSessionId(sessionId)
    if userId != "" {
	return userId, true 
    }
    return "", false
}
