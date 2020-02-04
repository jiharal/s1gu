package cmd

import "fmt"

func (n S1GU) createAPIUser(appPath string) string {
	body := `
	package api

	import (
		"context"
		"database/sql"
		"net/http"
		"strings"
	
		"github.com/gomodule/redigo/redis"
		"` + fmt.Sprintf("%s/model", appPath) + `"
		"` + fmt.Sprintf("%s/system", appPath) + `"
		"github.com/pkg/errors"
		uuid "github.com/satori/go.uuid"
	)
	
	type (
		UserModule struct {
			db    *sql.DB
			cache *redis.Pool
			name  string
		}
		UserDataParam struct {
			ID       uuid.UUID ` + fmt.Sprintf("`json:%s`", `"id"`) + `
			Name     string    ` + fmt.Sprintf("`json:%s`", `"name"`) + `
			Email    string    ` + fmt.Sprintf("`json:%s`", `"email"`) + `
			Password string    ` + fmt.Sprintf("`json:%s`", `"password"`) + `
			By       uuid.UUID ` + fmt.Sprintf("`json:%s`", `"by"`) + `
		}
	
		UserLoginParam struct {
			Email      string ` + fmt.Sprintf("`json:%s`", `"email"`) + `
			Password   string ` + fmt.Sprintf("`json:%s`", `"password"`) + `
			RememberMe bool   ` + fmt.Sprintf("`json:%s`", `"remember_me"`) + `
		}
	)
	
	func NewUserModule(db *sql.DB, cache *redis.Pool) *UserModule {
		return &UserModule{
			db:    db,
			cache: cache,
			name:  "module/user",
		}
	}
	
	func (m UserModule) Register(ctx context.Context, param UserDataParam) (model.UserModelResponse, *Error) {
		hash, err := HashPassword(ctx, param.Password)
		if err != nil {
			return model.UserModelResponse{}, NewErrorWrap(err, m.name, "create/hash",
				MessageGeneralError, http.StatusInternalServerError)
		}
		user := model.UserModel{
			Name:      param.Name,
			Email:     param.Email,
			Password:  hash,
			CreatedBy: system.DefaultID(),
		}
	
		resp, err := user.Insert(ctx, m.db)
		if err != nil {
			status := http.StatusInternalServerError
			message := MessageGeneralError
			if strings.Contains(err.Error(), "duplicate") {
				status = http.StatusConflict
				message = MessageAccountExists
			}
			return model.UserModelResponse{}, NewErrorWrap(err, m.name, "insert/customer",
				message, status)
		}
	
		return resp.Response(), nil
	}
	
	func (m UserModule) Login(ctx context.Context, param UserLoginParam) (model.UserModelResponse, string, *Error) {
		var resp model.UserModelResponse
		user, err := model.GetOneUserByEmail(ctx, m.db, param.Email)
		if err != nil {
			status := http.StatusInternalServerError
			message := MessageGeneralError
	
			if errors.Cause(err) == sql.ErrNoRows {
				status = http.StatusUnauthorized
				message = MessageInvalidLogin
			}
	
			return resp, "", NewErrorWrap(err, m.name, "login/user",
				message, status)
		}
		var token string
		err = ComparePassword(ctx, user.Password, param.Password)
		if err != nil {
			return resp, "", NewErrorWrap(err, m.name, "login/password",
				MessageInvalidLogin, http.StatusUnauthorized)
		}
		session := Session{
			User: user.Response(),
		}
	
		token, err = NewSession(ctx, session)
		if err != nil {
			return resp, "", NewErrorWrap(err, m.name, "login/session",
				MessageGeneralError, http.StatusInternalServerError)
		}
		return resp, token, nil
	}
	
	func (m UserModule) List(ctx context.Context) ([]model.UserModelResponse, *Error) {
		Users, err := model.GetAllUser(ctx, m.db)
		if err != nil {
			return nil, NewErrorWrap(err, m.name, "list/query",
				MessageGeneralError, http.StatusInternalServerError)
		}
	
		UserResponse := []model.UserModelResponse{}
	
		for _, User := range Users {
			UserResponse = append(UserResponse, User.Response())
		}
		return UserResponse, nil
	}
	
	func (m UserModule) Detail(ctx context.Context, id uuid.UUID) (model.UserModelResponse, *Error) {
		user, err := model.GetOneUser(ctx, m.db, id)
		if err != nil {
			status := http.StatusInternalServerError
			message := MessageGeneralError
	
			if err == sql.ErrNoRows {
				status = http.StatusNotFound
				message = http.StatusText(status)
			}
	
			return model.UserModelResponse{}, NewErrorWrap(err, m.name, "detail/query",
				message, status)
		}
	
		return user.Response(), nil
	}
	
	func (m UserModule) Create(ctx context.Context, param UserDataParam) (model.UserModelResponse, *Error) {
	
		hash, err := HashPassword(ctx, param.Password)
		if err != nil {
			return model.UserModelResponse{}, NewErrorWrap(err, m.name, "create/hash",
				MessageGeneralError, http.StatusInternalServerError)
		}
		user := model.UserModel{
			Name:      param.Name,
			Email:     param.Email,
			Password:  hash,
			CreatedBy: system.DefaultID(),
		}
	
		resp, err := user.Insert(ctx, m.db)
		if err != nil {
			return model.UserModelResponse{}, NewErrorWrap(err, m.name, "create",
				MessageGeneralError, http.StatusInternalServerError)
		}
	
		return resp.Response(), nil
	}
	
	func (m UserModule) Update(ctx context.Context, param UserDataParam) *Error {
		hash, err := HashPassword(ctx, param.Password)
		if err != nil {
			return NewErrorWrap(err, m.name, "create/hash",
				MessageGeneralError, http.StatusInternalServerError)
		}
		user, err := model.GetOneUser(ctx, m.db, param.ID)
		if err != nil {
			return NewErrorWrap(err, m.name, "create/hash",
				MessageGeneralError, http.StatusInternalServerError)
		}
	
		if param.Name != "" || param.Email != "" || param.Password != "" {
			user.Name = param.Name
			user.Email = param.Email
			user.Password = hash
			user.UpdatedBy = uuid.NullUUID{UUID: param.By, Valid: true}
		}
	
		err = user.Update(ctx, m.db)
		if err != nil {
			return NewErrorWrap(err, m.name, "update",
				MessageGeneralError, http.StatusInternalServerError)
		}
	
		return nil
	}
	
	func (m UserModule) Delete(ctx context.Context, id uuid.UUID) *Error {
		err := model.DeleteUser(ctx, m.db, id)
		if err != nil {
			return NewErrorWrap(err, m.name, "delete",
				MessageGeneralError, http.StatusInternalServerError)
		}
		return nil
	}
	`
	return body
}

func (n S1GU) createErrorFile() string {
	return `
	package api

	import (
		"github.com/pkg/errors"
	)
	
	type (
		Error struct {
			Err        error
			StatusCode int
			Message    string
		}
	)
	
	func (e *Error) Error() string {
		return e.Message
	}
	
	func NewError(err error, message string, status int) *Error {
		return &Error{
			Err:        err,
			Message:    message,
			StatusCode: status,
		}
	}
	
	func NewErrorWrap(err error, prefix, suffix, message string, status int) *Error {
		return &Error{
			Err:        errors.Wrapf(err, "%s/%s", prefix, suffix),
			Message:    message,
			StatusCode: status,
		}
	}
	`
}

func (n S1GU) createAPIInitFile() string {
	return `
	package api

	import (
		"database/sql"
	
		"github.com/KancioDevelopment/lib-angindai/logging"
		"github.com/gomodule/redigo/redis"
	)
	
	type (
		InitOption struct {
			SessionExpire string
		}
	)
	
	var (
		logger    *logging.Logger
		dbPool    *sql.DB
		cachePool *redis.Pool
		cfg       InitOption
	)
	
	func Init(lg *logging.Logger, db *sql.DB, cache *redis.Pool, opt InitOption) {
		logger = lg
		dbPool = db
		cachePool = cache
		cfg = opt
	}
	`
}

func (n S1GU) createAPIResponse() string {
	return `package api

	import (
		"context"
		"encoding/json"
		"net/http"
	
		"golang.org/x/crypto/bcrypt"
	)
	
	type (
		BaseResponse struct {
			Errors []string ` + fmt.Sprintf("`json:%s`", `"errors,omitempty"`) + `
		}
		Response struct {
			Status       string ` + fmt.Sprintf("`json:%s`", `"status,omitempty"`) + `
			BaseResponse ` + fmt.Sprintf("`json:%s`", `"errors,omitempty"`) + `
			Data         interface{} ` + fmt.Sprintf("`json:%s`", `"data"`) + `
		}
	)
	
	var (
		MessageGeneralError  = "Ada kesalahan, Silahakan coba beberapa saat lagi."
		MessageUnauthorized  = "Silahkan login terlebih dahulu atau login ulang."
		MessageInvalidLogin  = "Email atau password anda salah."
		MessageAccountExists = "Akun anda sudah terdaftar, silahkan login."
	)
	
	// RespondError writes / respond with JSON-formatted request of given message & http status.
	func RespondError(w http.ResponseWriter, message string, status int) {
		resp := Response{
			Status: http.StatusText(status),
			Data:   nil,
			BaseResponse: BaseResponse{
				Errors: []string{
					message,
				},
			},
		}
	
		w.WriteHeader(status)
		w.Header().Set("Content-Type", "application/json")
	
		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			logger.Err.WithError(err).Println("Encode response error.")
			return
		}
	}
	
	func ComparePassword(ctx context.Context, hash string, password string) error {
		return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	}
	`
}

func (n S1GU) createNewSessionAPIFile(appPath string) string {
	body := `package api

	import (
		"context"
		"database/sql"
		"encoding/json"
		"fmt"
		"time"
	
		"github.com/gomodule/redigo/redis"
		"` + fmt.Sprintf("%s/model", appPath) + `"
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
			User model.UserModelResponse ` + fmt.Sprintf("`json:%s`", `"user"`) + `
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
	`
	return body
}
