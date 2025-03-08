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
    sql := fmt.Sprintf("INSERT INTO users (%s, %s, %s) VALUES ($1, $2, $3);", DB_users_id, DB_users_email, DB_users_password)
    _, err := p.db.Exec(p.ctx, sql, user.Id, user.Email, user.Password)
    if err != nil {
	return "", fmt.Errorf("Couldn't create user: %w", err) 
    }
    return user.Id, nil
}

func (p *PostgreUsers) DeleteUser(userId string) error {
    sql := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", DB_users_name)
    _, err := p.db.Exec(p.ctx, sql, userId)
    if err != nil {
	err = fmt.Errorf("Couldn't delete user: %w", err)
    }
    return err
}

func (p *PostgreUsers) ReadUserById(id string) (*models.User, error) {
    sql := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", DB_users_name)
    row := p.db.QueryRow(p.ctx, sql, id)
    var user *models.User
    err := row.Scan(user)
    if err != nil {
	return nil, fmt.Errorf("Couldn't read user by Id: %w", err)
    }
    return user, nil
}

func (p *PostgreUsers) FindUser(user *models.User) (string, error) {
    sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s=$1 AND %s=$2;", DB_users_id, DB_users_name, DB_users_email, DB_users_password)
    row := p.db.QueryRow(p.ctx, sql, user.Email, user.Password) 
    var userId string
    err := row.Scan(&userId)
    if err != nil {
	return "", fmt.Errorf("Couldn't find user. %w", err)
    }
    return userId, nil 
}
