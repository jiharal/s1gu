
	package router

	import (
		"github.com/graphql-go/graphql"
	)

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
	