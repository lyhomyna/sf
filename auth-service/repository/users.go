package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lyhomyna/sf/auth-service/database"
	"github.com/lyhomyna/sf/auth-service/models"
)

type PostgreUsers struct {
    ctx context.Context
    db *pgx.Conn
}

func GetUsersDao(ctx context.Context) *PostgreUsers {
    dbConnection := database.ConnectToDb(ctx)
    if dbConnection == nil {
	return nil
    }
    
    return &PostgreUsers{ctx, dbConnection}
}

func (p *PostgreUsers) CreateUser(user *models.DbUser) (string, error) {
    _, err := p.db.Exec(p.ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3);", user.Id, user.Email, user.Password)
    if err != nil {
	return "", fmt.Errorf("Couldn't create user: %w", err) 
    }

    return user.Id, nil
}

func (p *PostgreUsers) DeleteUser(userId string) error {
    _, err := p.db.Exec(p.ctx, "DELETE FROM users WHERE id=$1;", userId)
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

func (p *PostgreUsers) FindUser(user *models.User) (string, error) {
    row := p.db.QueryRow(p.ctx, "SELECT * FROM users WHERE email=$1, pwd=$2", user.Email, user.Password) 
    
    var dbUser *models.DbUser
    err := row.Scan(dbUser)
    if err != nil {
	return "", fmt.Errorf("Couldn't find user. %w", err)
    }
    return dbUser.Id, nil
}
