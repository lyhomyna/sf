package service

import (
	"context"
	"net/http"

	"github.com/lyhomyna/sf/auth-service/models"
	"github.com/lyhomyna/sf/auth-service/repository"
	usersService "github.com/lyhomyna/sf/auth-service/service/user"
	sessionService "github.com/lyhomyna/sf/auth-service/service/session"
)
type Services struct {
    Users UserService
    Sessions SessionService
}

func GetServices(ctx context.Context, repos *repository.Repositories) *Services {
    return &Services {
	Users: usersService.NewUserService(repos.Users),
	Sessions: sessionService.NewSessionService(repos.Sessions),
    }
}

type UserService interface {
    CreateUser(req *http.Request) (string, *models.HTTPError)
    DeleteUser(userId string) *models.HTTPError   
    GetUserByEmail(email string) (*models.DbUser, *models.HTTPError)
    GetUserById(userId string) (*models.User, *models.HTTPError)
}

type SessionService interface {
    CreateSession(userId string) (string, *models.HTTPError)
    DeleteSession(userId string) *models.HTTPError
    IsSessionExists(sessionId string) (string, bool)
}
