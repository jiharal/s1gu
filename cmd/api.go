package cmd

import (
	"fmt"
	"log"
	"os"
	path "path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jiharal/s1gu/utils"
)

var apiCommand = &cobra.Command{
	Use:   "api [api name]",
	Short: "Create api aplication",
	Args:  cobra.MinimumNArgs(1),
	Run:   createAPI,
}

func createAPI(cmd *cobra.Command, args []string) {
	output := cmd.OutOrStderr()
	getEnv, _ := os.Getwd()

	pathModel := path.Join(getEnv, "/api")
	var appPath string

	gps := utils.GetGOPATHs()
	if len(gps) == 0 {
		log.Fatal("GOPATH environment variable is not set or empty")
	}

	for _, gpath := range gps {
		gsrcpath := path.Join(gpath, "src")
		if strings.HasPrefix(strings.ToLower(pathModel), strings.ToLower(gsrcpath)) {
			appPath = strings.Replace(getEnv[len(gsrcpath)+1:], string(path.Separator), "/", -1)
		}
	}

	if utils.IsExist(pathModel) {
		log.Print("Do you want to add it? [Yes|No] ")
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}

	fileName := "api." + strings.ToLower(args[0]) + ".go"

	replacer := strings.NewReplacer(
		"{{.AppPath}}", appPath,
		"{{.ApiName}}", strings.Title(args[0]),
		"{{.ApiNameLower}}", strings.ToLower(args[0]),
	)
	fileContent := `
	package api

	import (
		"context"
		"database/sql"
		"net/http"
	
		"github.com/gomodule/redigo/redis"
		uuid "github.com/satori/go.uuid"
	
		"{{.AppPath}}/model"
		"{{.AppPath}}/util"
	)

	type (
		// {{.ApiName}}Module is ...
		{{.ApiName}}Module struct {
			db    *sql.DB
			cache *redis.Pool
			name  string
		}

		// {{.ApiName}}DataParam is ...
		{{.ApiName}}DataParam struct {
			ID    uuid.UUID ` + fmt.Sprintf("`json:%s`", `"id"`) + `
			Name  string ` + fmt.Sprintf("`json:%s`", `"name"`) + `
			By    uuid.UUID ` + fmt.Sprintf("`json:%s`", `"by"`) + `
		}
	)

	// New{{.ApiName}}Module is ...
	func New{{.ApiName}}Module(db *sql.DB, cache *redis.Pool) *{{.ApiName}}Module {
		return &{{.ApiName}}Module{
			db:    db,
			cache: cache,
			name:  "module/{{.ApiNameLower}}",
		}
	}

	// List is a ...
	func (m {{.ApiName}}Module) List(ctx context.Context, filter model.FilterOption) ([]model.{{.ApiName}}ModelResponse, *Error) {
		{{.ApiNameLower}}s, err := model.GetAll{{.ApiName}}(ctx, m.db, filter)
		if err != nil {
			return nil, NewErrorWrap(err, m.name, "list/query", 
				MessageGeneralError, http.StatusInternalServerError)
		}

		{{.ApiNameLower}}Response := []model.{{.ApiName}}ModelResponse{}

		for _, {{.ApiNameLower}} := range {{.ApiNameLower}}s {
			{{.ApiNameLower}}Response = append({{.ApiNameLower}}Response, {{.ApiNameLower}}.Response())
		}

		return {{.ApiNameLower}}Response, nil
	}
	

	// Detail is ...
	func (m {{.ApiName}}Module) Detail(ctx context.Context, id uuid.UUID) (model.{{.ApiName}}ModelResponse, *Error) {
		{{.ApiNameLower}}, err := model.GetOne{{.ApiName}}(ctx, m.db, id)
		if err != nil {
			status := http.StatusInternalServerError
			message := MessageGeneralError

			if err == sql.ErrNoRows {
				status = http.StatusNotFound
				message = http.StatusText(status)
			}

			return model.{{.ApiName}}ModelResponse{}, NewErrorWrap(err, m.name, "detail/query",
				message, status)
		}

		return {{.ApiNameLower}}.Response(), nil
	}
	
	// Create is ...
	func (m {{.ApiName}}Module) Create(ctx context.Context, param {{.ApiName}}DataParam) (model.{{.ApiName}}ModelResponse, *Error) {
		
		{{.ApiNameLower}} := model.{{.ApiName}}Model{
			Name:      param.Name,
			CreatedBy:   param.By,
		}

		err := {{.ApiNameLower}}.Insert(ctx, m.db)
		if err != nil {
			return model.{{.ApiName}}ModelResponse{}, NewErrorWrap(err, m.name, "create",
				MessageGeneralError, http.StatusInternalServerError)
		}

		return {{.ApiNameLower}}.Response(), nil 
	}
	
	// Update is ...
	func (m {{.ApiName}}Module) Update(ctx context.Context, param {{.ApiName}}DataParam) *Error {
		{{.ApiNameLower}} := model.{{.ApiName}}Model{
			ID: param.ID,
			Name:      param.Name,
			UpdatedBy:   uuid.NullUUID{Valid: true, UUID: param.By},
		}

		err := {{.ApiNameLower}}.Update(ctx, m.db)
		if err != nil {
			return NewErrorWrap(err, m.name, "update",
			MessageGeneralError, http.StatusInternalServerError)
		}

		return nil
	}
	
	// Delete is ...
	func (m {{.ApiName}}Module) Delete(ctx context.Context, id uuid.UUID) *Error {
		err := model.Delete{{.ApiName}}ByID(ctx, m.db, id)
		if err != nil {
			return NewErrorWrap(err, m.name, "delete",
				MessageGeneralError, http.StatusInternalServerError)
		}
		return nil
	}
	`

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(pathModel, fileName), "\x1b[0m")
	utils.WriteToFile(path.Join(pathModel, fileName), replacer.Replace(fileContent))
}
