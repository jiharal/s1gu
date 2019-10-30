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

var newCommand = &cobra.Command{
	Use:   "new",
	Short: "Create new project",
	Args:  cobra.MinimumNArgs(1),
	Run:   createNewFramework,
}

func createNewFramework(cmd *cobra.Command, args []string) {

	// Read new file
	maingo := `
	package main

	import (
		"{{.Appname}}/cmd"
	)

	func main() {
		cmd.Execute()
	}`

	fileFileModelUser := `
	package model

	import (
		"context"
		"database/sql"
		"time"

		"github.com/lib/pq"
		"github.com/satori/go.uuid"
	)

	type (
		UserModel struct {
			ID        uuid.UUID     ` + fmt.Sprintf("`json:%s`", `"id"`) + `
			Name      string        ` + fmt.Sprintf("`json:%s`", `"name"`) + `
			CreatedAt time.Time     ` + fmt.Sprintf("`json:%s`", `"created_at"`) + `
			CreatedBy uuid.UUID     ` + fmt.Sprintf("`json:%s`", `"created_by"`) + `
			UpdatedAt pq.NullTime   ` + fmt.Sprintf("`json:%s`", `"updated_at"`) + `
			UpdatedBy uuid.NullUUID ` + fmt.Sprintf("`json:%s`", `"updated_by"`) + `
		}
	)

	func GetAllUser(ctx context.Context, db *sql.DB) ([]UserModel, error) {
		var userList []UserModel
		query := ` + fmt.Sprintf("`%s`", `SELECT id, name FROM "user"`) + `
		
		rows, err := db.QueryContext(ctx, query)
		if err != nil {
			return userList, err
		}
		defer rows.Close()

		for rows.Next() {
			var user UserModel
			err := rows.Scan(
				&user.ID,
				&user.Name,
			)
			if err != nil {
				return userList, err
			}
			userList = append(userList, user)
		}
		return userList, nil
	}

	func GetOneUser(ctx context.Context, db *sql.DB, ID uuid.UUID) (UserModel, error) {
		var user UserModel
		query := ` + fmt.Sprintf("`%s`", `SELECT id, name FROM "user" WHERE id=$1`) + `
		err := db.QueryRowContext(ctx, query, ID).Scan(
			&user.ID,
			&user.Name,
		)
		if err != nil {
			return user, err
		}
		return user, nil
	}

	func (usr UserModel) Insert(ctx context.Context, db *sql.DB) (uuid.UUID, error) {
		var id uuid.UUID

		query := ` + fmt.Sprintf("`%s`", `INSERT INTO "user"(name, created_by, created_at)VALUES($1, $2, now()) RETURNING id`) + `
		err := db.QueryRowContext(ctx, query,
			usr.Name,
			usr.CreatedBy).Scan(&id)
		if err != nil {
			return id, err
		}
		return id, nil
	}

	func (usr UserModel) Update(ctx context.Context, db *sql.DB) error {
		query := ` + fmt.Sprintf("`%s`", `UPDATE "user" SET(name, updated_by, updated_at)=($1, $2, now()) WHERE id=$3`) + `
		_, err := db.ExecContext(ctx, query,
			usr.Name,
			usr.UpdatedBy,
			usr.ID)
		if err != nil {
			return err
		}
		return nil
	}

	func DeleteUser(ctx context.Context, db *sql.DB, ID uuid.UUID) error {
		query := ` + fmt.Sprintf("`%s`", `DELETE FROM "user" WHERE id = $1`) + `
		_, err := db.ExecContext(ctx, query, ID)
		if err != nil {
			return err
		}
		return nil
	}
	`

	fileRouterUserRest := `
	package router

	import (
		"encoding/json"
		"net/http"

		"github.com/satori/go.uuid"

		"{{.AppPath}}/model"
		"github.com/gorilla/mux"
	)

	func GetAllUser(w http.ResponseWriter, r *http.Request) {
		users, err := model.GetAllUser(r.Context(), DbPool)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string][]model.UserModel{"user": users})
	}
	
	func GetOneUser(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := uuid.FromString(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		user, err := model.GetOneUser(r.Context(), DbPool, id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]model.UserModel{"user": user})
	}
	
	func InsertUser(w http.ResponseWriter, r *http.Request) {
		var user model.UserModel
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()
	
		id, err := user.Insert(r.Context(), DbPool)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		user.ID = id
	
		respondWithJSON(w, http.StatusCreated, map[string]model.UserModel{"user": user})
	}
	
	func UpdateUser(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := uuid.FromString(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	
		var user model.UserModel
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		defer r.Body.Close()
	
		user.ID = id
	
		if err := user.Update(r.Context(), DbPool); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	
		respondWithJSON(w, http.StatusOK, map[string]string{"user": "success"})
	}
	
	func DeleteUser(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := uuid.FromString(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	
		if err := model.DeleteUser(r.Context(), DbPool, id); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	
		respondWithJSON(w, http.StatusOK, map[string]string{"user": "success"})
	}
	
	`

	fileRouterUserHandler := `
	package router

	import (
		"{{.AppPath}}/model"
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
					Name:    p.Args["name"].(string),
					CreatedBy:   uuid.FromStringOrNil(p.Args["created_by"].(string)),
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
					ID:          uuid.FromStringOrNil(p.Args["id"].(string)),
					Name:        p.Args["full_name"].(string),
					UpdatedBy:   uuid.NullUUID{UUID: uuid.FromStringOrNil(p.Args["updated_by"].(string)), Valid: true},
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
	`

	fileRouterHandler := `
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
	`

	fileRouterInit := `
	package router

	import (
		"database/sql"
		"encoding/json"
		"net/http"

		"github.com/AuthScureDevelopment/lib-arjuna/logging"
		"github.com/gomodule/redigo/redis"
	)

	var (
		Logger    *logging.Logger
		DbPool    *sql.DB
		CachePool *redis.Pool
	)

	func Init(db *sql.DB, cachePool *redis.Pool, logger *logging.Logger) {
		DbPool = db
		CachePool = cachePool
		Logger = logger
	}

	type ErrorMethod struct {
		Errors interface{} ` + fmt.Sprintf("`json:%s`", `"errors"`) + `
	}

	type ResponseMethod struct {
		Data interface{} ` + fmt.Sprintf("`json:%s`", `"data"`) + `
	}

	func respondWithError(w http.ResponseWriter, code int, message interface{}) {
		response, _ := json.Marshal(ErrorMethod{Errors: message})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	}

	func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
		response, _ := json.Marshal(ResponseMethod{Data: payload})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(response)
	}
	`

	configFile := `
	[app]
	name = "{{.AppName}}"
	version = "0.0.1"
	port = 9100

	[database]
	host = "localhost"
	port = 26257
	username = "root"
	password = ""
	name = "{{.DbName}}_db"
	sslmode = "disable"

	[cache]
	host = "localhost"
	port = 6379
	password = ""
	max_idle = 100
	idle_timeout = 5
	enabled = true
	expire_time = 60
	idempotency_expiry = 86400
	`

	cmdFile := `
	package cmd

	import (
		"database/sql"
		"fmt"
		"net/http"
		"os"
	
		log "github.com/sirupsen/logrus"
	
		"{{.AppPath}}/router"
		"github.com/AuthScureDevelopment/lib-arjuna/cache"
		"github.com/AuthScureDevelopment/lib-arjuna/db"
		"github.com/AuthScureDevelopment/lib-arjuna/logging"
		"github.com/gomodule/redigo/redis"
		"github.com/gorilla/mux"
		"github.com/graphql-go/handler"
		homedir "github.com/mitchellh/go-homedir"
		"github.com/spf13/cobra"
		"github.com/spf13/viper"
	)
	
	var (
		dbPool    *sql.DB
		cfgFile   string
		cachePool *redis.Pool
		logger    *logging.Logger
	)
	
	var rootCmd = &cobra.Command{
		Use:   "{{.AppName}}",
		Short: "Simple golang app",
		Run: func(cmd *cobra.Command, args []string) {
			router.Init(dbPool, cachePool, logger)
			schema, err := router.InitSchema()
			if err != nil {
				log.Fatalln("Initiate schema error:", err)
			}
	
			r := mux.NewRouter()
	
			// GraphQL API
			graphqlHandler := handler.New(&handler.Config{
				Schema: &schema,
				Pretty: true,
			})
	
			r.Handle("/graphql", graphqlHandler)
			fs := http.FileServer(http.Dir("static"))
			r.Handle("/", fs)
	
			// RESTful API
			r.HandleFunc("/users", router.GetAllUser).Methods("GET")
			r.HandleFunc("/users/{id}", router.GetOneUser).Methods("GET")
			r.HandleFunc("/users", router.InsertUser).Methods("POST")
			r.HandleFunc("/users/{id}", router.UpdateUser).Methods("POST")
			r.HandleFunc("/users/{id}", router.DeleteUser).Methods("DELETE")
	
			fmt.Println("Listening on", fmt.Sprintf("http://localhost:%d", viper.GetInt("app.port")))
			http.ListenAndServe(fmt.Sprintf(":%d", viper.GetInt("app.port")), r)
		},
	}
	
	func Execute() {
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	
	func init() {
		cobra.OnInitialize(initConfig, initDB, initCache, initLogger)
		rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file(default is $HOME/.{{.AppName}}.config.toml)")
		rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	}
	
	func initConfig() {
		viper.SetConfigType("toml")
	
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			home, err := homedir.Dir()
			if err != nil {
				panic(err)
			}
			viper.AddConfigPath(".")
			viper.AddConfigPath(home)
			viper.SetConfigName(".{{.AppName}}")
		}
		viper.AutomaticEnv()
	
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("using config file: ", viper.ConfigFileUsed())
		}
	}
	
	func initDB() {
		dbOptions := db.DBOptions{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetInt("database.port"),
			Username: viper.GetString("database.username"),
			Password: viper.GetString("database.password"),
			DBName:   viper.GetString("database.name"),
			SSLMode:  viper.GetString("database.sslmode"),
		}
		dbConn, err := db.Connect(dbOptions)
		if err != nil {
			fmt.Println("Error conn to DB", err)
			panic(err)
		}
		dbPool = dbConn
	}
	
	func initCache() {
		cacheOptions := cache.CacheOptions{
			Host:        viper.GetString("cache.host"),
			Port:        viper.GetInt("cache.port"),
			Password:    viper.GetString("cache.password"),
			MaxIdle:     viper.GetInt("cache.max_idle"),
			IdleTimeout: viper.GetInt("cache.idle_timeout"),
			Enabled:     viper.GetBool("cache.enabled"),
		}
		cachePool = cache.Connect(cacheOptions)
	}
	
	func initLogger() {
		logger = logging.New()
		logger.Out.Formatter = new(log.JSONFormatter)
		logger.Err.Formatter = new(log.JSONFormatter)
	}
	`

	output := cmd.OutOrStderr()
	apppath, packpath, err := utils.CheckEnv(args[0])
	if err != nil {
		log.Fatalf("%s", err)
	}

	if utils.IsExist(apppath) {
		log.Print("Do you want to add it? [Yes|No] ")
		if !utils.AskForConfirmation() {
			os.Exit(2)
		}
	}

	fileConfigName := "." + strings.ToLower(args[0]) + ".toml"

	cmdReplaceContent := strings.NewReplacer(
		"{{.AppPath}}", packpath,
		"{{.AppName}}", strings.ToLower(args[0]),
	)

	configReplaceContent := strings.NewReplacer(
		"{{.AppName}}", args[0],
		"{{.DbName}}", strings.ToLower(args[0]),
	)

	// Create root APP
	os.MkdirAll(apppath, 0755)

	// Create cmd directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "cmd"), 0755)

	// Create file cmd.go in cmd directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "cmd", "cmd.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "cmd", "cmd.go"), cmdReplaceContent.Replace(string(cmdFile)))

	// Create model directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "model"), 0755)

	// Create file model.user.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "model", "model.user.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "model", "model.user.go"), string(fileFileModelUser))

	// Create route directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "router"), 0755)

	// Create file init.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "router", "init.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "router", "init.go"), fileRouterInit)

	// Create file handler.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "router", "handler.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "router", "handler.go"), string(fileRouterHandler))

	// Create file graphql.user.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "router", "graphql.user.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "router", "graphql.user.go"), strings.Replace(string(fileRouterUserHandler), "{{.AppPath}}", packpath, -1))

	// Create file restful.user.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "router", "restful.user.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "router", "restful.user.go"), strings.Replace(string(fileRouterUserRest), "{{.AppPath}}", packpath, -1))

	// Create file main.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "main.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "main.go"), strings.Replace(string(maingo), "{{.Appname}}", packpath, -1))

	// Create file config.toml
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, fileConfigName), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, fileConfigName), configReplaceContent.Replace(string(configFile)))

	fmt.Println("New application successfully created!")
}
