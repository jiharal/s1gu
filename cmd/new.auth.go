package cmd

import "fmt"

func (n S1GU) createInitAuth() string {
	return `package auth

	import (
		"database/sql"
	
		"github.com/KancioDevelopment/lib-angindai/logging"
		"github.com/gomodule/redigo/redis"
	)
	
	type (
		InitOption struct{}
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

func (n S1GU) createAuthMiddleware(appPath string) string {
	return `package auth

	import (
		"net/http"
	
		"` + fmt.Sprintf("%s/api", appPath) + `"
		"github.com/pkg/errors"
	)
	
	func AuthenticationMiddleware(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			token := r.Header.Get("Authorization")
			var (
				session api.Session
				err     error
			)
	
			if token != "" {
				session, err = api.GetSession(ctx, token)
			} else {
				err = api.ErrUnauthorized
			}
	
			if err != nil {
				if errors.Cause(err) == api.ErrUnauthorized {
					api.RespondError(w, api.MessageUnauthorized, http.StatusInternalServerError)
					return
				}
			}
	
			ctx = api.SetRequesterContext(ctx, session)
			ctx = api.SetContextToken(ctx, token)
			r = r.WithContext(ctx)
	
			next.ServeHTTP(w, r)
		})
	}
	`
}
