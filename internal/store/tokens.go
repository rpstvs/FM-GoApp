package store

import (
	"database/sql"
	"time"

	"github.com/rpstvs/fm-goapp/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userId int, scope string) error
}

func (t *PostgresTokenStore) CreateNewToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userId, ttl, scope)

	if err != nil {
		return nil, err
	}

	err = t.Insert(token)
	return token, err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES($1,$2,$3,$4)
	`
	_, err := t.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userId int, scope string) error {
	query := `
	DELETE from tokens
	WHERE scope =$1 AND user_id = $2
	`

	_, err := t.db.Exec(query, scope, userId)

	return err
}
