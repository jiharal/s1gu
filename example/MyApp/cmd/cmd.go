package cmd

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/AuthScureDevelopment/authscure-go/example/MyApp/router"
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
	Use:   "myapp",
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

		// Access API
		r.HandleFunc("/access", router.GetAllAccess).Methods("GET")
		r.HandleFunc("/access/{id}", router.GetOneAccess).Methods("GET")
		r.HandleFunc("/access", router.InsertAccess).Methods("POST")
		r.HandleFunc("/access/{id}", router.UpdateAccess).Methods("POST")
		r.HandleFunc("/access/{id}", router.DeleteAccess).Methods("DELETE")

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
