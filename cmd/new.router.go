package cmd

import "fmt"

func CreateInitFile(appPath string) string {
	initFile :=
		`
		package router

		import (
			"context"
			"database/sql"
			"encoding/json"
			"io/ioutil"
			"net/http"

			"github.com/KancioDevelopment/lib-angindai/logging"
			"github.com/asaskevich/govalidator"
			"github.com/gomodule/redigo/redis"
			"` + fmt.Sprintf("%s/api", appPath) + `"
			"github.com/pkg/errors"
		)

		type (
			InitOption struct{}
		)

		var (
			logger    *logging.Logger
			dbPool    *sql.DB
			cachePool *redis.Pool
			cfg       InitOption

			userService *api.UserModule
		)

		func Init(lg *logging.Logger, db *sql.DB, cache *redis.Pool, opt InitOption) {
			logger = lg
			dbPool = db
			cachePool = cache
			cfg = opt

			userService = api.NewUserModule(dbPool, cachePool)
		}

		// ParseBodyData parse json-formatted request body into given struct.
		func ParseBodyData(ctx context.Context, r *http.Request, data interface{}) error {
			bBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return errors.Wrap(err, "read")
			}

			err = json.Unmarshal(bBody, data)
			if err != nil {
				return errors.Wrap(err, "json")
			}

			valid, err := govalidator.ValidateStruct(data)
			if err != nil {
				return errors.Wrap(err, "validate")
			}

			if !valid {
				return errors.New("invalid data")
			}

			return nil
		}
	`

	return initFile
}

func CreateHandlerFile(appPath string) string {
	handlerFile :=
		`
		package router

		import (
			"encoding/json"
			"net/http"

			"github.com/graphql-go/graphql"
			"` + fmt.Sprintf("%s/api", appPath) + `"
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
		`
	return handlerFile
}

func CreateGraphQLFile(appPath, appName string) string {
	graphQLFile := `
	package router

	import (
		"github.com/graphql-go/graphql"
		"` + fmt.Sprintf("%s/api", appPath) + `"
		uuid "github.com/satori/go.uuid"
	)
	
	var (
		userType = graphql.NewObject(graphql.ObjectConfig{
			Name: "User",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.String,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"created_at": &graphql.Field{
					Type: graphql.String,
				},
				"created_by": &graphql.Field{
					Type: graphql.String,
				},
				"updated_at": &graphql.Field{
					Type: graphql.String,
				},
				"updated_by": &graphql.Field{
					Type: graphql.String,
				},
			},
		})
	
		userListField = &graphql.Field{
			Type: graphql.NewList(userType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context
				data, err := userService.List(ctx)
				if err != nil {
					return data, err
				}
				return data, nil
			},
		}
	
		userDetailField = &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context
				id, err := uuid.FromString(p.Args["id"].(string))
				if err != nil {
					return nil, err
				}
	
				data, err := userService.Detail(ctx, id)
				if err != nil {
					return data, err
				}
	
				return data, nil
			},
		}
	
		userCreateField = &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"created_by": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context
				user := api.UserDataParam{
					Name: p.Args["name"].(string),
					By:   uuid.FromStringOrNil(p.Args["created_by"].(string)),
				}
	
				data, err := userService.Create(ctx, user)
				if err != nil {
					return data, err
				}
				return data, nil
			},
		}
	
		userUpdateField = &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"updated_by": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context
	
				user := api.UserDataParam{
					ID:   uuid.FromStringOrNil(p.Args["id"].(string)),
					Name: p.Args["name"].(string),
					By:   uuid.FromStringOrNil(p.Args["updated_by"].(string)),
				}
	
				err := userService.Update(ctx, user)
				if err != nil {
					return nil, err
				}
				return nil, nil
			},
		}
	
		userDeleteField = &graphql.Field{
			Type: userType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctx := p.Context
				id, err := uuid.FromString(p.Args["id"].(string))
				if err != nil {
					return nil, err
				}
	
				err = userService.Delete(ctx, id)
				if err != nil {
					return nil, err
				}
				return nil, nil
			},
		}
	)	
	`
	return graphQLFile
}

func CreateRestFile(appPath string) string {
	body := `
	package router

	import (
		"net/http"

		"github.com/gorilla/mux"
		"` + fmt.Sprintf("%s/api", appPath) + `"
		"github.com/pkg/errors"
		uuid "github.com/satori/go.uuid"
	)

	func HandlerUserLogin(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
		ctx := r.Context()

		var param api.UserLoginParam

		err := ParseBodyData(ctx, r, &param)
		if err != nil {
			return nil, api.NewError(errors.Wrap(err, "vehicle/create/param"),
				http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		_, token, aErr := userService.Login(ctx, param)
		if aErr != nil {
			return nil, aErr
		}

		data := struct {
			AccessToken string ` + fmt.Sprintf("`json:%s`", `"access_token"`) + `
		}{
			AccessToken: token,
		}

		return data, nil
	}
	func HandlerUserRegister(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
		ctx := r.Context()

		var param api.UserDataParam

		err := ParseBodyData(ctx, r, &param)
		if err != nil {
			return nil, api.NewError(errors.Wrap(err, "vehicle/create/param"),
				http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		return userService.Register(ctx, param)
	}

	func HandlerUserList(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
		ctx := r.Context()
		return userService.List(ctx)
	}

	func HandlerUserDetail(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
		ctx := r.Context()
		vars := mux.Vars(r)

		id, err := uuid.FromString(vars["id"])
		if err != nil {
			return nil, api.NewError(errors.Wrap(err, "vehicle/detail"),
				http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		return userService.Detail(ctx, id)
	}

	func HandlerUserCreate(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
		ctx := r.Context()

		var param api.UserDataParam

		err := ParseBodyData(ctx, r, &param)
		if err != nil {
			return nil, api.NewError(errors.Wrap(err, "vehicle/create/param"),
				http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		param.By = api.GetContextRequesterID(ctx)
		return userService.Create(ctx, param)
	}

	func HandlerUserUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
		ctx := r.Context()
		vars := mux.Vars(r)

		id, err := uuid.FromString(vars["id"])
		if err != nil {
			return nil, api.NewError(errors.Wrap(err, "vehicle/update"),
				http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
		var param api.UserDataParam
		err = ParseBodyData(ctx, r, &param)
		if err != nil {
			return nil, api.NewError(errors.Wrap(err, "vehicle/update/param"),
				http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		param.ID = id
		param.By = api.GetContextRequesterID(ctx)
		return nil, userService.Update(ctx, param)
	}

	func HandlerUserDelete(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
		ctx := r.Context()
		vars := mux.Vars(r)

		id, err := uuid.FromString(vars["id"])
		if err != nil {
			return nil, api.NewError(errors.Wrap(err, "vehicle/delete"),
				http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}

		return nil, userService.Delete(ctx, id)
	}`
	return body
}
