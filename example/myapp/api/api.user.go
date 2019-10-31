package api

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/jiharal/s1gu/example/myapp/model"
	"github.com/jiharal/s1gu/example/myapp/system"
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
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Email    string    `json:"email"`
		Password string    `json:"password"`
		By       uuid.UUID `json:"by"`
	}

	UserLoginParam struct {
		Email      string `json:"email"`
		Password   string `json:"password"`
		RememberMe bool   `json:"remember_me"`
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
