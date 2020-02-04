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
	appName := args[0]
	// Create root APP
	os.MkdirAll(apppath, 0755)

	// Create api directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "api"), 0755)
	// Create file init.go in api directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "init.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "api", "init.go"), CreateAPIInitFile())
	// Create file init.go in api directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "api.user.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "api", "api.user.go"), CreateAPIUser(packpath))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "error.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "api", "error.go"), CreateErrorFile())
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "response.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "api", "response.go"), CreateAPIResponse())
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "session.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "api", "session.go"), CreateNewSessionAPIFile(packpath))

	// Create cmd directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "cmd"), 0755)
	// Create file cmd.go in cmd directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "cmd", "cmd.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "cmd", "cmd.go"), CreateNewCMD(packpath, strings.ToLower(args[0])))

	// Create auth directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "auth"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "init.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "auth", "init.go"), CreateInitAuth())
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "middleware.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "auth", "middleware.go"), CreateAuthMiddleware(packpath))

	// Create model directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "model"), 0755)

	// Create file model.user.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "model", "model.user.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "model", "model.user.go"), CreateNewModelFile())

	// Create route directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "router"), 0755)

	// Create file init.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "router", "init.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "router", "init.go"), CreateInitFile(packpath))

	// Create file handler.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "router", "handler.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "router", "handler.go"), CreateHandlerFile(packpath))

	// Create file graphql.user.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "router", "graphql.user.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "router", "graphql.user.go"), CreateGraphQLFile(packpath, appName))

	// Create file restful.user.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "router", "handler.user.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "router", "handler.user.go"), CreateRestFile(packpath))

	// Create system directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "system"), 0755)
	// Create file id.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "system", "id.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "system", "id.go"), CreateGeneralID())

	// Create files directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "files"), 0755)
	// Create file restful.user.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "files", "db.sql"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "files", "db.sql"), CreateDBFile())

	// Create file main.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "main.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "main.go"), CreateMainPage(packpath))

	// Create file config.toml
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, fileConfigName), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, fileConfigName), CreateConfigFile(appName))

	fmt.Println("New application successfully created!")
}
