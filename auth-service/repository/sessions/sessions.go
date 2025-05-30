package sessions

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lyhomyna/sf/auth-service/database"
	"github.com/lyhomyna/sf/auth-service/models"
)

var CookieSeessionIdName = "session-id"

type PostgreSessions struct {
    ctx context.Context
    pool *pgxpool.Pool
}

func GetSessionsDao(ctx context.Context) *PostgreSessions {
    connPool := database.ConnectToDb(ctx)
    if connPool == nil {
	return nil
    }
    
    return &PostgreSessions{ctx, connPool}
}

func (p *PostgreSessions) CreateSession(userId string) (string, error) {
    sessionId := uuid.NewString()
    expiresAt := time.Now().Add(24 * time.Hour)

    sql := fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES ($1, $2, $3)", 
	database.DB_sessions_name, 
	database.DB_sessions_id, 
	database.DB_sessions_user_id, 
	database.DB_sessions_expires_at)

    _, err := p.pool.Exec(p.ctx, sql, sessionId, userId, expiresAt)

    if err != nil {
	return "", fmt.Errorf("Couldn't create session: %w", err)
    }

    return sessionId, nil
}

func (p *PostgreSessions) DeleteSession(sessionId string) error {
    sql := fmt.Sprintf("DELETE FROM %s WHERE id=$1", database.DB_sessions_name)

    execRes, err := p.pool.Exec(p.ctx, sql, sessionId)

    if err != nil {
	err = fmt.Errorf("Couldn't delete session: %w", err) 
    } else if execRes.RowsAffected() != 1 {
	err = errors.New("Nothing was deleted from sessions table")
    }

    return err
}

// TODO: ...
func (p *PostgreSessions) UserIdFromSessionId(sessionId string) (string, error) {
    sql := fmt.Sprintf("SELECT %s, %s, %s FROM %s WHERE id=$1", 
	database.DB_sessions_id,
	database.DB_sessions_user_id, 
	database.DB_sessions_expires_at,
	database.DB_sessions_name)
    row := p.pool.QueryRow(p.ctx, sql, sessionId)

    var	session models.Session 
    err := row.Scan(&session.Id, &session.UserId, &session.ExpiresAt)
    if err != nil {
	if errors.Is(err, pgx.ErrNoRows) {
            return "", fmt.Errorf("Session not found")
        }
        return "", fmt.Errorf("DB error: %w", err)
    }
    
    if session.ExpiresAt.Before(time.Now()) {
        go func() {
            _, _ = p.pool.Exec(context.Background(), `DELETE FROM sessions WHERE id = $1`, sessionId)
        }()
        return "", fmt.Errorf("Session expired")
    }

    return session.UserId, nil
} 
