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

var modelCommand = &cobra.Command{
	Use:   "model [model name]",
	Short: "Create model application",
	Args:  cobra.MinimumNArgs(1),
	Run:   CreateModel,
}

func CreateModel(cmd *cobra.Command, args []string) {
	output := cmd.OutOrStderr()
	getEnv, _ := os.Getwd()

	pathModel := path.Join(getEnv, "/model")

	if utils.IsExist(pathModel) {
		log.Print("Do you want to add it? [Yes|No] ")
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}

	filename := "model." + strings.ToLower(args[0]) + ".go"
	replacer := strings.NewReplacer(
		"{{.ModelName}}", strings.Title(args[0]),
		"{{.ModelNameLower}}", strings.ToLower(args[0]),
	)

	var fileContent = `package model
		import (
			"context"
			"database/sql"
			"fmt"
			"time"

			"github.com/lib/pq"
			"github.com/pkg/errors"
			uuid "github.com/satori/go.uuid"
		)

		type (

			{{.ModelName}}Model struct{
				ID        uuid.UUID
				Name      string     
				CreatedAt time.Time
				CreatedBy uuid.UUID
				UpdatedAt pq.NullTime
				UpdatedBy uuid.NullUUID
			}

			{{.ModelName}}ModelResponse struct{
				ID        uuid.UUID     ` + fmt.Sprintf("`json:%s`", `"id"`) + `
				Name      string        ` + fmt.Sprintf("`json:%s`", `"name"`) + `
				CreatedAt time.Time     ` + fmt.Sprintf("`json:%s`", `"created_at"`) + `
				CreatedBy uuid.UUID     ` + fmt.Sprintf("`json:%s`", `"created_by"`) + `
				UpdatedAt time.Time   	` + fmt.Sprintf("`json:%s`", `"updated_at"`) + `
				UpdatedBy uuid.UUID 		` + fmt.Sprintf("`json:%s`", `"updated_by"`) + `
			}
		)


		// Convert {{.ModelNameLower}} model into json-friendly formatted response struct (without null data type).
		func ({{.ModelNameLower}} *{{.ModelName}}Model) Response() {{.ModelName}}ModelResponse {
			return {{.ModelName}}ModelResponse{
				ID:        {{.ModelNameLower}}.ID,
				Name:      {{.ModelNameLower}}.Name,
				CreatedAt: {{.ModelNameLower}}.CreatedAt,
				CreatedBy: {{.ModelNameLower}}.CreatedBy,
				UpdatedAt: {{.ModelNameLower}}.UpdatedAt.Time,
				UpdatedBy: {{.ModelNameLower}}.UpdatedBy.UUID,
			}
		}

		// Implements api/SessionData interface.
		func (am {{.ModelName}}Model) Identifier() uuid.UUID {
			return am.ID
		}

		func GetAll{{.ModelName}}(ctx context.Context, db *sql.DB, filter FilterOption) ([]{{.ModelName}}Model, error) {
			if filter.Dir != "ASC" && filter.Dir != "DESC" {
				return nil, errors.New("Invalid order by parameter")
			}
		
			query := fmt.Sprintf(` + fmt.Sprintf("`%s`", `SELECT
				id,
				name,
				created_at,
				created_by,
				updated_at,
				updated_by
			FROM
				{{.ModelNameLower}}
			WHERE
				id = CASE WHEN $1 <> '' THEN $1 ELSE id END
			ORDER BY
				id %s
			LIMIT $2 OFFSET $3`) + `, filter.Dir)
		
			rows, err := db.QueryContext(ctx, query, filter.Search, filter.Limit, filter.Offset)
			if err != nil {
				return nil, errors.Wrap(err, "model/{{.ModelNameLower}}/list")
			}
			defer rows.Close()
		
			var {{.ModelNameLower}}s []{{.ModelName}}Model
		
			for rows.Next() {
				var {{.ModelNameLower}} {{.ModelName}}Model
		
				err = rows.Scan(
					&{{.ModelNameLower}}.ID,
					&{{.ModelNameLower}}.Name,
					&{{.ModelNameLower}}.CreatedAt,
					&{{.ModelNameLower}}.CreatedBy,
					&{{.ModelNameLower}}.UpdatedAt,
					&{{.ModelNameLower}}.UpdatedBy,
				)
				if err != nil {
					return nil, errors.Wrap(err, "model/{{.ModelNameLower}}/list/scan")
				}
		
				{{.ModelNameLower}}s = append({{.ModelNameLower}}s, {{.ModelNameLower}})
			}
			return {{.ModelNameLower}}s, nil
		}

		func GetOne{{.ModelName}}(ctx context.Context, db *sql.DB, id uuid.UUID) ({{.ModelName}}Model, error) {
			query := ` + fmt.Sprintf("`%s`", `
				SELECT
					id,
					name,
					created_at,
					created_by,
					updated_at,
					updated_by
				FROM
					{{.ModelNameLower}}
				WHERE
					id = $1`) + `
		
			var {{.ModelNameLower}} {{.ModelName}}Model
		
			err := db.QueryRowContext(ctx, query, id).Scan(
				&{{.ModelNameLower}}.ID,
				&{{.ModelNameLower}}.Name,
				&{{.ModelNameLower}}.CreatedAt,
				&{{.ModelNameLower}}.CreatedBy,
				&{{.ModelNameLower}}.UpdatedAt,
				&{{.ModelNameLower}}.UpdatedBy,
			)
			if err != nil {
				return {{.ModelName}}Model{}, errors.Wrap(err, "model/{{.ModelNameLower}}/query/id")
			}
		
			return {{.ModelNameLower}}, nil
		}

		func ({{.ModelNameLower}} *{{.ModelName}}Model) Insert(ctx context.Context, db *sql.DB) error {
			query := ` + fmt.Sprintf("`%s`", `
			INSERT INTO {{.ModelNameLower}} (
				name,
				created_by,
				created_at
			) VALUES (
				$1, $2, now()
			) RETURNING
				id,
				created_at`) + `
		
			err := db.QueryRowContext(ctx, query,
				{{.ModelNameLower}}.Name,
			).Scan(
				&{{.ModelNameLower}}.ID,
				&{{.ModelNameLower}}.CreatedAt,
			)
			if err != nil {
				return errors.Wrap(err, "model/{{.ModelNameLower}}/insert")
			}
		
			return nil
		}

		func ({{.ModelNameLower}} *{{.ModelName}}Model) Update(ctx context.Context, db *sql.DB) error {
			query := ` + fmt.Sprintf("`%s`", `
				UPDATE
					{{.ModelNameLower}}
				SET
					name = $1,
					updated_by = $2,
					updated_at = NOW()
				WHERE
					id=$3`) + `
		
			_, err := db.ExecContext(ctx, query,
				{{.ModelNameLower}}.Name,
				{{.ModelNameLower}}.UpdatedBy,
				{{.ModelNameLower}}.ID,
			)
			if err != nil {
				return errors.Wrap(err, "model/{{.ModelNameLower}}/update")
			}
		
			return nil
		}


		func Delete{{.ModelName}}ByID(ctx context.Context, db *sql.DB, id uuid.UUID) error {
			query := "DELETE FROM {{.ModelNameLower}} WHERE id = $1"
		
			_, err := db.ExecContext(ctx, query, id)
			if err != nil {
				return errors.Wrap(err, "model/{{.ModelNameLower}}/delete")
			}
		
			return nil
		}`

	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(pathModel, filename), "\x1b[0m")
	utils.WriteToFile(path.Join(pathModel, filename), replacer.Replace(fileContent))
}
