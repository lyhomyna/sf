package repository

import (
	"context"
	
	"github.com/lyhomyna/sf/auth-service/models"
)

func GetSiglog() *models.Siglog {
    ctx := context.Background()
    usersDao := GetUsersDao(ctx)
    sessionsDao := GetSessionsDao(ctx)

    if usersDao == nil || sessionsDao == nil {
	return nil
    }

    siglog := &models.Siglog{
	Users: usersDao,
	Sessions: sessionsDao,
    }

    return siglog
}
