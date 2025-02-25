package database

import (
	"context"

	"github.com/lyhomyna/sf/auth-service/database/models"
	"github.com/lyhomyna/sf/auth-service/database/postgres"
)

func GetDao() *models.Siglog {
    ctx := context.Background()
    usersDao := postgres.GetUsersDao(ctx)
    sessionsDao := postgres.GetSessionsDao(ctx)

    siglog := &models.Siglog{
	Users: usersDao,
	Sessions: sessionsDao,
    }

    return siglog
}
