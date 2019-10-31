package router

import (
	"encoding/json"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/jiharal/s1gu/example/myapp/api"
)

type (
	HandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, *api.Error)
)

func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var errs []string

	// Ignore error from form parsing as it's insignificant.
	r.ParseForm()

	data, err := fn(w, r)
	if err != nil {
		logger.Err.WithError(err.Err).Println("Serve error.")
		errs = append(errs, err.Error())
		w.WriteHeader(err.StatusCode)
		resp := api.Response{
			Status: http.StatusText(err.StatusCode),
			Data:   data,
			BaseResponse: api.BaseResponse{
				Errors: errs,
			},
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			logger.Err.WithError(err).Println("Encode response error.")
			return
		}
	} else {
		resp := api.Response{
			Status: http.StatusText(200),
			Data:   data,
			BaseResponse: api.BaseResponse{
				Errors: errs,
			},
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(&resp); err != nil {
			logger.Err.WithError(err).Println("Encode response error.")
			return
		}
	}
}

func InitSchema() (graphql.Schema, error) {
	queryFields := graphql.Fields{
		"user_list":   userListField,
		"user_detail": userDetailField,
	}

	mutationFields := graphql.Fields{
		"user_create": userCreateField,
		"user_update": userUpdateField,
		"user_delete": userDeleteField,
	}

	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Query",
			Fields: queryFields,
		},
	)

	mutationType := graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: mutationFields,
		},
	)

	return graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)
}
