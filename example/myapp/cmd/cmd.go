package cmd

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/KancioDevelopment/lib-angindai/cache"
	"github.com/KancioDevelopment/lib-angindai/db"
	"github.com/KancioDevelopment/lib-angindai/logging"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"github.com/jiharal/s1gu/example/myapp/api"
	"github.com/jiharal/s1gu/example/myapp/auth"
	"github.com/jiharal/s1gu/example/myapp/router"
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
	Use:   "myapp",
	Short: "Simple golang app",
	PreRun: func(cmd *cobra.Command, args []string) {
		api.Init(logger, dbPool, cachePool, api.InitOption{
			SessionExpire: viper.GetString("login.expire_time"),
		})
		router.Init(logger, dbPool, cachePool, router.InitOption{})
	},
	Run: func(cmd *cobra.Command, args []string) {

		schema, err := router.InitSchema()
		if err != nil {
			log.Fatalln("Initiate schema error:", err)
		}

		r := mux.NewRouter()

		r.Handle("/auth/login", router.HandlerFunc(router.HandlerUserLogin)).Methods("POST")
		r.Handle("/auth/signup", router.HandlerFunc(router.HandlerUserRegister)).Methods("POST")
		// Subrouter for API paths.
		apiRouter := r.PathPrefix("/api/v1").Subrouter()

		// Middlewares.
		apiRouter.Use(auth.AuthenticationMiddleware)

		// RESTful API
		apiRouter.Handle("/users", router.HandlerFunc(router.HandlerUserList)).Methods("GET")
		apiRouter.Handle("/users", router.HandlerFunc(router.HandlerUserCreate)).Methods("POST")
		apiRouter.Handle("/users/{id}", router.HandlerFunc(router.HandlerUserDetail)).Methods("GET")
		apiRouter.Handle("/users/{id}", router.HandlerFunc(router.HandlerUserUpdate)).Methods("POST")
		apiRouter.Handle("/users/{id}", router.HandlerFunc(router.HandlerUserDelete)).Methods("DELETE")

		// GraphQL API
		graphqlHandler := handler.New(&handler.Config{
			Schema: &schema,
			Pretty: true,
		})

		r.Handle("/graphql", graphqlHandler)
		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static"))))

		fmt.Println("Listening on", fmt.Sprintf("http://:%d", viper.GetInt("app.port")))
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
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file(default is $HOME/.myapp.config.toml)")
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
		viper.SetConfigName(".myapp")
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
