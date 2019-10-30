package router

import (
	"github.com/AuthScureDevelopment/authscure-go/example/MyApp/model"
	"github.com/graphql-go/graphql"
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

			userList, err := model.GetAllUser(ctx, DbPool)
			if err != nil {
				return nil, err
			}
			return userList, nil
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

			user, err := model.GetOneUser(ctx, DbPool, id)
			if err != nil {
				return nil, err
			}
			return user, nil
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
			user := model.UserModel{
				Name:      p.Args["name"].(string),
				CreatedBy: uuid.FromStringOrNil(p.Args["created_by"].(string)),
			}
			id, err := user.Insert(ctx, DbPool)
			if err != nil {
				return nil, err
			}
			return model.UserModel{ID: id}, nil
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

			user := model.UserModel{
				ID:        uuid.FromStringOrNil(p.Args["id"].(string)),
				Name:      p.Args["full_name"].(string),
				UpdatedBy: uuid.NullUUID{UUID: uuid.FromStringOrNil(p.Args["updated_by"].(string)), Valid: true},
			}

			err := user.Update(ctx, DbPool)
			if err != nil {
				return nil, err
			}
			return user, nil
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

			err = model.DeleteUser(ctx, DbPool, id)
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
	}
)
