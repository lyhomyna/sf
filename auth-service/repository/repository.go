package repository

import (
	"context"

	"github.com/lyhomyna/sf/auth-service/models"
	"github.com/lyhomyna/sf/auth-service/repository/sessions"
	"github.com/lyhomyna/sf/auth-service/repository/users"
)

type Repositories struct {
    Users UserDao 
    Sessions SessionDao 
}

func GetRepos(ctx context.Context) *Repositories {
    usersDao := users.GetUsersDao(ctx)
    sessionsDao := sessions.GetSessionsDao(ctx)

    if usersDao == nil || sessionsDao == nil {
	return nil
    }

    return &Repositories{
	Users: usersDao,
	Sessions: sessionsDao,
    }
}

type UserDao interface {
    CreateUser(user *models.DbUser) (string, error)
    DeleteUser(userId string) error
    ReadUserById(id string) (*models.User, error)
    FindUser(user *models.User) (string, error)
    GetUserByEmail(email string) (*models.DbUser, error)
}

type SessionDao interface {
    CreateSession(userId string) (string, error)
    DeleteSession(sessionId string) error
    UserIdFromSessionId(sessionId string) (string, error)
}
