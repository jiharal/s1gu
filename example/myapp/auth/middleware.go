package auth

import (
	"net/http"

	"github.com/jiharal/s1gu/example/myapp/api"
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
