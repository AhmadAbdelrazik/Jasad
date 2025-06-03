package model

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Session struct {
	UserID int
	Role   Role
}

// For Redis client to be able to marshal it.
func (s Session) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

type TokenRepository struct {
	redis *redis.Client
}

// GenerateToken returns session token for the given user.
func (r *TokenRepository) GenerateToken(user *User) (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	token := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes)

	hash := sha256.Sum256([]byte(token))

	session := &Session{
		UserID: user.ID,
		Role:   user.Role,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.redis.Set(ctx, string(hash[:]), session, 3*24*time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("failed to set value on redis: %w", err)
	}

	return token, nil
}

func (r *TokenRepository) GetSessionFromToken(token string) (*Session, error) {
	hash := sha256.Sum256([]byte(token))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionStr, err := r.redis.Get(ctx, string(hash[:])).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, ErrNotFound
		default:
			return nil, fmt.Errorf("failed to get value on redis: %w", err)
		}
	}

	var session Session
	if err := json.Unmarshal([]byte(sessionStr), &session); err != nil {
		return nil, fmt.Errorf("failed to marshal session: %w", err)
	}

	return &session, nil
}
