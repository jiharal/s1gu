
	package api

	import (
		"context"
		"database/sql"
		"net/http"
	
		"golang.org/x/crypto/bcrypt"
	
		"github.com/gomodule/redigo/redis"
		"github.com/pkg/errors"
	
		"github.com/AuthScureDevelopment/authscure-go/example/MyApp/model"
		"github.com/AuthScureDevelopment/authscure-go/example/MyApp/util"
	)

	type (
		HelloModule struct {
			db    *sql.DB
			cache *redis.Pool
		}
	
		HelloParam struct {
			ID    string `json:"id"`
		}
	)

	
	func NewHelloModule(db *sql.DB, cache *redis.Pool) *HelloModule {
		return &HelloModule{
			db:    db,
			cache: cache,
		}
	}

	// Get all Hello
	func (m HelloModule) GetAllHello(ctx context.Context, param HelloParam) ([]model.HelloModelResponse, *Error) {

		return []model.HelloModelResponse{}, nil
	}
	

	func (m HelloModule) GetOneHello(ctx context.Context, param HelloParam) (model.HelloModelResponse, *Error) {
	
		return model.HelloModelResponse{}, nil
	}
	
	func (m HelloModule) Insert(ctx context.Context, param HelloParam) (uuid.UUID, *Error) {
		var id uuid.UUID
		return id, nil
	}
	
	func (m HelloModule) Update(ctx context.Context, param HelloParam) *Error {
	
		return nil
	}
	
	func (m HelloModule) Delete(ctx context.Context, param HelloParam) *Error {
	
		return nil
	}
	