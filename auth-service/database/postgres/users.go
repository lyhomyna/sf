package postgres

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/lyhomyna/sf/auth-service/database/models"
)

type PostgreUsers struct {
    ctx context.Context
    db *pgx.Conn
}

func GetUsersDao(ctx context.Context) *PostgreUsers {
    dbConnection := connectToDb(ctx)
    if dbConnection == nil {
	return nil
    }
    
    return &PostgreUsers{ctx, dbConnection}
}

func (p *PostgreUsers) CreateUser(user *models.User) (string, error) {
    _, err := p.db.Exec(p.ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3);", user.Id, user.Email, user.Password)
    if err != nil {
	return "", fmt.Errorf("Couldn't create user: %w", err) 
    }

    return user.Id, nil
}

func (p *PostgreUsers) DeleteUser(user *models.User) error {
    _, err := p.db.Exec(p.ctx, "DELETE FROM users WHERE id=$1;", user.Id)
    if err != nil {
	err = fmt.Errorf("Couldn't delete user: %w", err)
    }
    return err
}

func (p *PostgreUsers) ReadUserById(id string) (*models.User, error) {
    row := p.db.QueryRow(p.ctx, "SELECT * FROM users WHERE id=$1;", id)

    var user *models.User
    err := row.Scan(user)
    if err != nil {
	return nil, fmt.Errorf("Couldn't read user by Id: %w", err)
    }
    
    return user, nil
}
