package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgreSessions struct {
    db *pgx.Conn
    ctx context.Context
}

func (p *PostgreSessions) CreateSession(userId string) (string, error) {
    sessionId := uuid.NewString()
    _, err := p.db.Exec(p.ctx, "INSERT INTO sessions (id, user_id) VALUES ($1, $2);", sessionId, userId)

    if err != nil {
	return "", fmt.Errorf("Couldn't create session: %w", err)
    }

    return sessionId, nil
}

func (p *PostgreSessions) DeleteSession(sessionId string) error {
    execRes, err := p.db.Exec(p.ctx, "DELETE FROM sessions WHERE id=$1;", sessionId)

    if err != nil {
	err = fmt.Errorf("Couldn't delete session: %w", err) 
    } else if execRes.RowsAffected() != 1 {
	err = errors.New("Nothing was deleted from sessions table.")
    }

    return err
}

func (p *PostgreSessions) UserIdFromSessionId(sessionId string) (string, error) {
    row := p.db.QueryRow(p.ctx, "SELECT user_id FROM sessions WHERE id=$1", sessionId)

    var	userId string

    err := row.Scan(&userId)
    if err != nil {
	return "", fmt.Errorf("Couldn't get username from session '%s': %w", sessionId, err)
    }

    return userId, nil
} 
