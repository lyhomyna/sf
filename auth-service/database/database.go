package database

import (
	"context"

	"github.com/lyhomyna/sf/auth-service/database/models"
	"github.com/lyhomyna/sf/auth-service/database/postgres"
)

func GetSiglog() *models.Siglog {
    ctx := context.Background()
    usersDao := postgres.GetUsersDao(ctx)
    sessionsDao := postgres.GetSessionsDao(ctx)

    if usersDao == nil || sessionsDao == nil {
	return nil
    }

    siglog := &models.Siglog{
	Users: usersDao,
	Sessions: sessionsDao,
    }

    return siglog
}
