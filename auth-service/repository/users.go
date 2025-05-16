package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lyhomyna/sf/auth-service/database"
	"github.com/lyhomyna/sf/auth-service/models"
)

type PostgreUsers struct {
    ctx context.Context
    pool *pgxpool.Pool
}

func GetUsersDao(ctx context.Context) *PostgreUsers {
    connPool := database.ConnectToDb(ctx)
    if connPool == nil {
	return nil
    }
    return &PostgreUsers{ctx, connPool}
}

func (p *PostgreUsers) CreateUser(user *models.DbUser) (string, error) {
    sql := fmt.Sprintf("INSERT INTO users (%s, %s, %s) VALUES ($1, $2, $3)", DB_users_id, DB_users_email, DB_users_password)
    _, err := p.pool.Exec(p.ctx, sql, user.Id, user.Email, user.Password)
    if err != nil {
	return "", fmt.Errorf("Couldn't create user: %w", err) 
    }
    return user.Id, nil
}

func (p *PostgreUsers) DeleteUser(userId string) error {
    sql := fmt.Sprintf("DELETE FROM %s WHERE id=$1", DB_users_name)
    _, err := p.pool.Exec(p.ctx, sql, userId)
    if err != nil {
	err = fmt.Errorf("Couldn't delete user: %w", err)
    }
    return err
}

func (p *PostgreUsers) ReadUserById(id string) (*models.User, error) {
    sql := fmt.Sprintf("SELECT %s, %s, %s FROM %s WHERE id=$1", DB_users_email, DB_users_password, DB_users_image_url, DB_users_name)
    row := p.pool.QueryRow(p.ctx, sql, id)
    var user models.User

    err := row.Scan(&user.Email, &user.Password, &user.ImageUrl)
    if err != nil {
	return nil, fmt.Errorf("Couldn't read user by Id: %w", err)
    }

    return &user, nil
}

func (p *PostgreUsers) FindUser(user *models.User) (string, error) {
    sql := fmt.Sprintf("SELECT %s FROM %s WHERE %s=$1 AND %s=$2", DB_users_id, DB_users_name, DB_users_email, DB_users_password)
    row := p.pool.QueryRow(p.ctx, sql, user.Email, user.Password) 
    var userId string
    err := row.Scan(&userId)
    if err != nil {
	return "", fmt.Errorf("Couldn't find user. %w", err)
    }
    return userId, nil 
}

func (p *PostgreUsers) GetUserByEmail(email string) (*models.DbUser, error) {
    sql := fmt.Sprintf("SELECT * FROM %s WHERE %s=$1", DB_users_name, DB_users_email)
    row := p.pool.QueryRow(p.ctx, sql, email)

    var dbUser models.DbUser
    err := row.Scan(&dbUser.Id, &dbUser.Email, &dbUser.Password, &dbUser.CreatedAt, &dbUser.ImageUrl)
    if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
	    return nil, nil    
	}
	return nil, fmt.Errorf("Couldn't get user. %w", err)
    }

    return &dbUser, nil
}






