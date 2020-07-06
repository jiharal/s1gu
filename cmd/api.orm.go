package cmd

import (
	"fmt"
	"log"
	"os"
	path "path/filepath"
	"strings"

	"github.com/jiharal/s1gu/utils"
	"github.com/spf13/cobra"
)

var apiORMCMD = &cobra.Command{
	Use:   "api-orm [api name]",
	Short: "It's depend on orm Model, this command is used to create API level bussiness.",
	Args:  cobra.MinimumNArgs(1),
	Run:   createORMAPI,
}

func createORMAPI(cmd *cobra.Command, args []string) {
	output := cmd.OutOrStderr()
	getEnv, _ := os.Getwd()

	pathModel := path.Join(getEnv, "/api")
	var appPath string
	gps := utils.GetGOPATHs()
	if len(gps) == 0 {
		log.Fatal("GOPATH env varialble is not set, please set GOPATH first")
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

		"{{.AppPath}}/model"
		"{{.AppPath}}/util"

		"github.com/go-pg/pg/v10"
		uuid "github.com/satori/go.uuid"
		"github.com/valyala/fasthttp"
	)

	type (
		// {{.ApiName}}Module is ...
		{{.ApiName}}Module struct {
			db   *pg.DB
			name string
		}

		// {{.ApiName}}DataParam is ...
		{{.ApiName}}DataParam struct {
			ID          uuid.UUID ` + fmt.Sprintf("`%s`", `json:"id"`) + `
			Name    		string    ` + fmt.Sprintf("`%s`", `json:"name"`) + `
			IsActive    bool    	` + fmt.Sprintf("`%s`", `json:"is_active"`) + `
			By          uuid.UUID ` + fmt.Sprintf("`%s`", `json:"by"`) + `
		}
	)

	// New{{.ApiName}}Module is ...
	func New{{.ApiName}}Module(db *pg.DB) *{{.ApiName}}Module {
		return &{{.ApiName}}Module{
			db:   db,
			name: "module/{{.ApiNameLower}}",
		}
	}

	// List is ...
	func (m {{.ApiName}}Module) List(ctx context.Context, filter model.FilterOption) (interface{}, *Error) {
		{{.ApiNameLower}}s, count, err := model.GetAll{{.ApiName}}(ctx, m.db, filter)
		if err != nil {
			return nil, NewErrorWrap(err, m.name, err.Error(), fasthttp.StatusInternalServerError)
		}
		{{.ApiNameLower}}Response := []model.{{.ApiName}}ModelResponse{}
		for _, {{.ApiNameLower}} := range {{.ApiNameLower}}s {
			{{.ApiNameLower}}Response = append({{.ApiNameLower}}Response, {{.ApiNameLower}}.Response())
		}
		return util.Paginate(filter.Page, filter.Limit, count, {{.ApiNameLower}}Response), nil
	}

	// Detail is ...
	func (m {{.ApiName}}Module) Detail(ctx context.Context, id string) (model.{{.ApiName}}ModelResponse, *Error) {
		{{.ApiNameLower}}, err := model.GetOne{{.ApiName}}(ctx, m.db, id)
		if err != nil {
			return model.{{.ApiName}}ModelResponse{}, NewErrorWrap(err, m.name, err.Error(), fasthttp.StatusInternalServerError)
		}

		return {{.ApiNameLower}}.Response(), nil
	}

	// Create is ...
	func (m {{.ApiName}}Module) Create(ctx context.Context, param {{.ApiName}}DataParam) (model.{{.ApiName}}ModelResponse, *Error) {
		{{.ApiNameLower}} := model.{{.ApiName}}Model{
			Name:     param.Name,
			IsActive:  param.IsActive,
			CreatedBy: param.By,
		}

		_, err := {{.ApiNameLower}}.Insert(ctx, m.db)
		if err != nil {
			return model.{{.ApiName}}ModelResponse{}, NewErrorWrap(err, m.name, err.Error(), fasthttp.StatusInternalServerError)
		}
		return {{.ApiNameLower}}.Response(), nil
	}

	// Update is ...
	func (m {{.ApiName}}Module) Update(ctx context.Context, param {{.ApiName}}DataParam) *Error {
		{{.ApiNameLower}}, err := model.GetOne{{.ApiName}}(ctx, m.db, param.ID.String())
		if err != nil {
			return NewErrorWrap(err, m.name, err.Error(), fasthttp.StatusInternalServerError)
		}
		{{.ApiNameLower}}.Name = param.Name
		{{.ApiNameLower}}.IsActive = param.IsActive
		{{.ApiNameLower}}.UpdatedBy =param.By

		err = {{.ApiNameLower}}.Update(ctx, m.db)
		if err != nil {
			return NewErrorWrap(err, m.name, err.Error(), fasthttp.StatusInternalServerError)
		}
		return nil
	}

	// Delete is ...
	func (m {{.ApiName}}Module) Delete(ctx context.Context, id string) *Error {
		_, err := model.{{.ApiName}}Delete(ctx, m.db, id)
		if err != nil {
			return NewErrorWrap(err, m.name, err.Error(), fasthttp.StatusInternalServerError)
		}
		return nil
	}
	`
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(pathModel, fileName), "\x1b[0m")
	utils.WriteToFile(path.Join(pathModel, fileName), replacer.Replace(fileContent))
}
