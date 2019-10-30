package router

import (
	"github.com/AuthScureDevelopment/authscure-go/example/MyApp/model"
	"github.com/graphql-go/graphql"
	uuid "github.com/satori/go.uuid"
)

var (
	accessType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Access",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			// Add another field here
		},
	})

	accessListField = &graphql.Field{
		Type: graphql.NewList(accessType),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context

			accessList, err := model.GetAllAccess(ctx, DbPool)
			if err != nil {
				return nil, err
			}
			return accessList, nil
		},
	}

	accessDetailField = &graphql.Field{
		Type: accessType,
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
			access, err := model.GetOneAccess(ctx, DbPool, id)
			if err != nil {
				return nil, err
			}
			return access, nil
		},
	}

	accessCreateField = &graphql.Field{
		Type: accessType,
		Args: graphql.FieldConfigArgument{
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			// Add another argument here
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context
			access := model.AccessModel{
				Name: p.Args["name"].(string),
			}
			id, err := access.Insert(ctx, DbPool)
			if err != nil {
				return nil, err
			}
			return model.AccessModel{ID: id}, nil
		},
	}

	accessUpdateField = &graphql.Field{
		Type: accessType,
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"name": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			// add another argument here
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			ctx := p.Context
			access := model.AccessModel{
				ID:   uuid.FromStringOrNil(p.Args["id"].(string)),
				Name: p.Args["name"].(string),
			}
			err := access.Update(ctx, DbPool)
			if err != nil {
				return nil, err
			}
			return access, nil
		},
	}

	accessDeleteField = &graphql.Field{
		Type: accessType,
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

			err = model.DeleteAccess(ctx, DbPool, id)
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
	}
)
