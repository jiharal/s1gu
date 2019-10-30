package cmd

import (
	"fmt"
	"log"
	"os"
	path "path/filepath"
	"strings"

	"github.com/AuthScureDevelopment/authscure-go/utils"
	"github.com/spf13/cobra"
)

var routerCommand = &cobra.Command{
	Use:   "router [router name] [graphql or rest]",
	Short: "Create router application",
	Args:  cobra.MinimumNArgs(2),
	Run:   CreateRouter,
}

func CreateRouter(cmd *cobra.Command, args []string) {
	output := cmd.OutOrStderr()
	getEnv, _ := os.Getwd()

	pathModel := path.Join(getEnv, "/router")
	var pathApp string

	gps := utils.GetGOPATHs()
	if len(gps) == 0 {
		log.Fatal("GOPATH environment variable is not set or empty")
	}

	for _, gpath := range gps {
		gsrcpath := path.Join(gpath, "src")
		if strings.HasPrefix(strings.ToLower(pathModel), strings.ToLower(gsrcpath)) {
			pathApp = strings.Replace(getEnv[len(gsrcpath)+1:], string(path.Separator), "/", -1)
		}
	}

	if utils.IsExist(pathModel) {
		log.Print("Do you want to add it? [Yes|No] ")
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}

	var filename, fileContent string
	var replacer *strings.Replacer

	if args[1] == "graphql" {
		filename = "graphql." + args[0] + ".go"
		typeName := args[0] + "Type"
		typeName1 := strings.Title(args[0])
		listName := args[0] + "ListField"
		detailName := args[0] + "DetailField"
		createName := args[0] + "CreateField"
		updateName := args[0] + "UpdateField"
		deleteName := args[0] + "DeleteField"
		dbVariable := strings.ToLower(args[0])
		modelName := strings.Title(args[0]) + "Model"
		getAll := "GetAll" + strings.Title(args[0])
		getOne := "GetOne" + strings.Title(args[0])
		delOne := "Delete" + strings.Title(args[0])
		replacer = strings.NewReplacer(
			"{{.Appname}}", pathApp,
			"{{.modelName}}", modelName,
			"{{.dbVariable}}", dbVariable,
			"{{.getAll}}", getAll,
			"{{.getOne}}", getOne,
			"{{.delOne}}", delOne,
			"{{.typeName}}", typeName,
			"{{.typeName1}}", typeName1,
			"{{.listName}}", listName,
			"{{.detailName}}", detailName,
			"{{.createName}}", createName,
			"{{.updateName}}", updateName,
			"{{.deleteName}}", deleteName,
		)

		fileContent = `
		package router
			import (
				"{{.Appname}}/model"
				"github.com/graphql-go/graphql"
				uuid "github.com/satori/go.uuid"
			)
			
			var (
				{{.typeName}} = graphql.NewObject(graphql.ObjectConfig{
					Name: "{{.typeName1}}",
					Fields: graphql.Fields{
						"id": &graphql.Field{
							Type: graphql.String,
						},
						// Add another field here
					},
				})
			
				{{.listName}} = &graphql.Field{
					Type: graphql.NewList({{.typeName}}),
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						ctx := p.Context

						{{.dbVariable}}List, err := model.{{.getAll}}(ctx, DbPool)
						if err != nil {
							return nil, err
						}
						return {{.dbVariable}}List, nil
					},
				}
			
				{{.detailName}} = &graphql.Field{
					Type: {{.typeName}},
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
						{{.dbVariable}}, err := model.{{.getOne}}(ctx, DbPool, id)
						if err != nil {
							return nil, err
						}
						return {{.dbVariable}}, nil
					},
				}
			
				{{.createName}} = &graphql.Field{
					Type: {{.typeName}},
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.NewNonNull(graphql.String),
						},
						// Add another argument here
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						ctx := p.Context
						{{.dbVariable}} := model.{{.modelName}}{
							Name:      p.Args["name"].(string),
						}
						id, err := {{.dbVariable}}.Insert(ctx, DbPool)
						if err != nil {
							return nil, err
						}
						return model.{{.modelName}}{ID: id}, nil
					},
				}
			
				{{.updateName}} = &graphql.Field{
					Type: {{.typeName}},
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
						{{.dbVariable}} := model.{{.modelName}}{
							ID:        uuid.FromStringOrNil(p.Args["id"].(string)),
							Name:      p.Args["name"].(string),
						}
						err := {{.dbVariable}}.Update(ctx, DbPool)
						if err != nil {
							return nil, err
						}
						return {{.dbVariable}}, nil
					},
				}
			
				{{.deleteName}} = &graphql.Field{
					Type: {{.typeName}},
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

						err = model.{{.delOne}}(ctx, DbPool, id)
						if err != nil {
							return nil, err
						}
						return nil, nil
					},
				}
			)
		`
	} else if args[1] == "rest" {
		filename = "handler." + strings.ToLower(args[0]) + ".go"
		replacer = strings.NewReplacer(
			"{{.AppPath}}", pathApp,
			"{{.routerName}}", strings.Title(args[0]),
			"{{.routerNameLower}}", strings.ToLower(args[0]),
		)
		fileContent = `
		package router

		import (
			"net/http"
			"strconv"

			"github.com/gorilla/mux"
			"github.com/pkg/errors"
			uuid "github.com/satori/go.uuid"
		
			"{{.AppPath}}/api"
			"{{.AppPath}}/util"
		)
		
		func Handler{{.routerName}}List(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
			ctx := r.Context()
			filter := ParseFilterFromForm(ctx, r.Form)
			return {{.routerNameLower}}Service.List(ctx, filter)
		}
		
		func Handler{{.routerName}}Detail(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
			ctx := r.Context()
			vars := mux.Vars(r)

			id, err := uuid.FromString(vars["id"])
			if err != nil {
				return nil, api.NewError(errors.Wrap(err, "{{.routerNameLower}}/detail"),
					http.StatusText(http.StatusNotFound), http.StatusNotFound)
			}

			return {{.routerNameLower}}Service.Detail(ctx, id)
		}
		
		func Handler{{.routerName}}Create(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
			ctx := r.Context()

			var param api.{{.routerName}}DataParam

			err := ParseBodyData(ctx, r, &param)
			if err != nil {
				return nil, api.NewError(errors.Wrap(err, "{{.routerNameLower}}/create/param"),
					http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}

			param.By = api.GetContextRequesterID(ctx)
			return {{.routerNameLower}}Service.Create(ctx, param)
		}
		
		func Handler{{.routerName}}Update(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
			ctx := r.Context()
			vars := mux.Vars(r)

			id, err := uuid.FromString(vars["id"])
			if err != nil {
				return nil, api.NewError(errors.Wrap(err, "{{.routerNameLower}}/update"),
					http.StatusText(http.StatusNotFound), http.StatusNotFound)
			}

			var param api.{{.routerName}}DataParam
			err = ParseBodyData(ctx, r, &param)
			if err != nil {
				return nil, api.NewError(errors.Wrap(err, "{{.routerNameLower}}/update/param"),
					http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}

			param.ID = id
			param.By = api.GetContextRequesterID(ctx)
			return nil, {{.routerNameLower}}Service.Update(ctx, param)
		}
		
		func Handler{{.routerName}}Delete(w http.ResponseWriter, r *http.Request) (interface{}, *api.Error) {
			ctx := r.Context()
			vars := mux.Vars(r)
		
			id, err := uuid.FromString(vars["id"])
			if err != nil {
				return nil, api.NewError(errors.Wrap(err, "{{.routerNameLower}}/delete"),
					http.StatusText(http.StatusNotFound), http.StatusNotFound)
			}
		
			return nil, {{.routerNameLower}}Service.Delete(ctx, id)
		}	`

	}
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(pathModel, filename), "\x1b[0m")
	utils.WriteToFile(path.Join(pathModel, filename), replacer.Replace(fileContent))
}
