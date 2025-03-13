package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lyhomyna/sf/auth-service/database"
)

var CookieSeessionIdName = "session-id"

type PostgreSessions struct {
    ctx context.Context
    db *pgx.Conn
}

func GetSessionsDao(ctx context.Context) *PostgreSessions {
    dbConnection := database.ConnectToDb(ctx)
    if dbConnection == nil {
	return nil
    }
    
    return &PostgreSessions{ctx, dbConnection}
}

func (p *PostgreSessions) CreateSession(userId string) (string, error) {
    sql := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES ($1, $2);", DB_sessions_name, DB_sessions_id, DB_sessions_userId)
    sessionId := uuid.NewString()
    _, err := p.db.Exec(p.ctx, sql, sessionId, userId)

    if err != nil {
	return "", fmt.Errorf("Couldn't create session: %w", err)
    }

    return sessionId, nil
}

func (p *PostgreSessions) DeleteSession(sessionId string) error {
    sql := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", DB_sessions_name)
    execRes, err := p.db.Exec(p.ctx, sql, sessionId)

    if err != nil {
	err = fmt.Errorf("Couldn't delete session: %w", err) 
    } else if execRes.RowsAffected() != 1 {
	err = errors.New("Nothing was deleted from sessions table")
    }

    return err
}

func (p *PostgreSessions) UserIdFromSessionId(sessionId string) (string, error) {
    sql := fmt.Sprintf("SELECT %s FROM %s WHERE id=$1", DB_sessions_userId, DB_sessions_name)
    row := p.db.QueryRow(p.ctx, sql, sessionId)

    var	userId string

    err := row.Scan(&userId)
    if err != nil {
	return "", fmt.Errorf("Couldn't get username from session '%s': %w", sessionId, err)
    }

    return userId, nil
} 
