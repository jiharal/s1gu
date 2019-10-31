package router

import (
	"github.com/graphql-go/graphql"
	"github.com/jiharal/s1gu/example/myapp/api"
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
