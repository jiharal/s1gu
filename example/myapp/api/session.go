package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jiharal/s1gu/example/myapp/model"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	SessionModule struct {
		db    *sql.DB
		cache *redis.Pool
		name  string
	}

	SessionData interface {
		// Identifier gets the ID field of current logged in user type.
		Identifier() uuid.UUID
	}
	Session struct {
		User model.UserModelResponse `json:"user"`
	}
)

var (
	ErrUnauthorized = errors.New("unauthorized")
)

const (
	contextSession     string = "session"
	contextRequesterID string = "requester_id"
	contextToken       string = "token"
)

// Get current session data's identifier / id value.
func (s *Session) RequesterID() uuid.UUID {
	return s.User.Identifier()
}

func SetRequesterContext(ctx context.Context, session Session) context.Context {
	ctx = SetContextSession(ctx, session)
	ctx = SetContextRequesterID(ctx, session.RequesterID())
	return ctx
}

func SetContextSession(ctx context.Context, session Session) context.Context {
	return context.WithValue(ctx, contextSession, session)
}
func SetContextRequesterID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, contextRequesterID, id)
}

func SetContextToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, contextToken, token)
}

func GetContextRequesterID(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(contextRequesterID).(uuid.UUID)
	return id
}

// SessionKey return session cache's key with token parameter.
func SessionKey(token string) string {
	return fmt.Sprintf("myapp:auth:session:%s", token)
}

// Create new session of given data.
func NewSession(ctx context.Context, s Session) (string, error) {
	token, err := HashPassword(ctx, string(time.Now().UnixNano()))
	if err != nil {
		return "", errors.Wrap(err, "session/conn")
	}

	conn, err := cachePool.GetContext(ctx)
	if err != nil {
		return "", errors.Wrap(err, "session/conn")
	}
	defer conn.Close()

	bSession, err := json.Marshal(s)
	if err != nil {
		return "", errors.Wrap(err, "session/marshal")
	}

	_, err = conn.Do("SETEX", SessionKey(token), cfg.SessionExpire, string(bSession))
	if err != nil {
		return "", errors.Wrap(err, "session/set")
	}

	return token, nil
}

func HashPassword(ctx context.Context, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "hash")
	}

	return string(hash), nil
}

// GetSession returns session detail of given token string.
func GetSession(ctx context.Context, token string) (Session, error) {
	conn, err := cachePool.GetContext(ctx)
	if err != nil {
		return Session{}, errors.Wrap(err, "session/conn")
	}
	defer conn.Close()
	// Check token availability session in redis.
	bSession, err := redis.Bytes(conn.Do("GET", SessionKey(token)))
	if err != nil {
		if err == redis.ErrNil {
			return Session{}, ErrUnauthorized
		}
		return Session{}, errors.Wrap(err, "session/get")
	}

	var session Session

	err = json.Unmarshal(bSession, &session)
	if err != nil {
		return Session{}, errors.Wrap(err, "session/unpack")
	}

	return session, nil
}
