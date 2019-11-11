package cmd

import "fmt"

func CreateConfigFile(appName string) string {
	configFile := `
	[app]
	name = "` + fmt.Sprintf("%s", appName) + `"
	version = "0.0.1"
	port = 9100
	
	[login]
	expire_time = "160000"

	[database]
	host = "localhost"
	port = 5432
	username = "postgres"
	password = ""
	name = "` + fmt.Sprintf("%s", appName) + `_db"
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
	return configFile
}
