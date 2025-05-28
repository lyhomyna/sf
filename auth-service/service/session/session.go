package session

import (
	"log"
	"net/http"

	"github.com/lyhomyna/sf/auth-service/models"
	"github.com/lyhomyna/sf/auth-service/repository"
)

type SessionService struct {
    dao repository.SessionDao 
}

func NewSessionService(sessionDao repository.SessionDao) *SessionService {
    return &SessionService{
	dao: sessionDao,
    }
}

func (s *SessionService) CreateSession(userId string) (string, *models.HTTPError) {
    sessionId, err := s.dao.CreateSession(userId)
    if err != nil {
	log.Println("Couldn't create session.", err.Error())
	return "", &models.HTTPError{
	    Code: http.StatusInternalServerError,
	    Message: err.Error(),
	}
    }
    log.Println("Session created")
    return sessionId, nil
}

func (s *SessionService) DeleteSession(sessionId string) *models.HTTPError {
    err := s.dao.DeleteSession(sessionId)
    if err != nil {
	log.Println("Couldn't delete session.", err.Error())
	return &models.HTTPError{
	    Code: http.StatusInternalServerError,
	    Message: err.Error(),
	}
    }
    return nil
}

func (s *SessionService) IsSessionExists(sessionId string) (string, bool) {
    userId, _ := s.dao.UserIdFromSessionId(sessionId)
    if userId != "" {
	return userId, true 
    }
    return "", false
}
