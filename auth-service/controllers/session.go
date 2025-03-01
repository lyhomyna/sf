package controllers

import (
	"net/http"

	"github.com/lyhomyna/sf/auth-service/database/models"
)

type SessionsController struct {
    Siglog *models.Siglog
}

func (this *SessionsController) Create(userId string, w http.ResponseWriter) error {
    panic("Not yet implemented.")
}

func (this *SessionsController) Delete(sessionId string, w http.ResponseWriter) error {
    panic("Not yet implemented.")
}

func (this *SessionsController) HasCookie(req *http.Request) bool {
    panic("Not yet implemented.")
}

func (this *SessionsController) GetUserId(sessionId string) (string, error) {
    panic("Not yet implemented.")   
}
