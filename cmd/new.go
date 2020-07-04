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

// S1GU is ...
type S1GU struct{}

var newCommand = &cobra.Command{
	Use:   "new",
	Short: "Create new project",
	Args:  cobra.MinimumNArgs(1),
	Run:   createNewFramework,
}

func createNewFramework(cmd *cobra.Command, args []string) {
	var newSigu S1GU

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
	utils.WriteToFile(path.Join(apppath, "api", "init.go"), newSigu.createAPIInitFile())
	// Create file init.go in api directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "api.user.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "api", "api.user.go"), newSigu.createAPIUser(packpath))
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "error.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "api", "error.go"), newSigu.createErrorFile())
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "response.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "api", "response.go"), newSigu.createAPIResponse())
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "session.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "api", "session.go"), newSigu.createNewSessionAPIFile(packpath))

	// Create cmd directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "cmd"), 0755)
	// Create file cmd.go in cmd directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "cmd", "cmd.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "cmd", "cmd.go"), newSigu.createNewCMD(packpath, strings.ToLower(args[0])))

	// Create auth directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "auth"), 0755)
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "init.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "auth", "init.go"), newSigu.createInitAuth())
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "api", "middleware.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "auth", "middleware.go"), newSigu.createAuthMiddleware(packpath))

	// Create model directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "model"), 0755)

	// Create file model.user.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "model", "model.user.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "model", "model.user.go"), newSigu.createNewModelFile())

	// Create route directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "router"), 0755)

	// Create file init.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "router", "init.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "router", "init.go"), newSigu.createInitFile(packpath))

	// Create file handler.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "router", "handler.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "router", "handler.go"), newSigu.createHandlerFile(packpath))

	// Create file graphql.user.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "router", "graphql.user.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "router", "graphql.user.go"), newSigu.createGraphQLFile(packpath, appName))

	// Create file restful.user.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "router", "handler.user.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "router", "handler.user.go"), newSigu.createRestFile(packpath))

	// Create system directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "system"), 0755)
	// Create file id.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "system", "id.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "system", "id.go"), newSigu.createGeneralID())
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "system", "validate.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "system", "validate.go"), newSigu.createValidate())

	// Create files directory
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", apppath+string(path.Separator), "\x1b[0m")
	os.Mkdir(path.Join(apppath, "files"), 0755)
	// Create file restful.user.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "files", "db.sql"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "files", "db.sql"), newSigu.createDBFile())

	// Create file main.go
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, "main.go"), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, "main.go"), newSigu.createMainPage(packpath))

	// Create file config.toml
	fmt.Fprintf(output, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", path.Join(apppath, fileConfigName), "\x1b[0m")
	utils.WriteToFile(path.Join(apppath, fileConfigName), newSigu.createConfigFile(appName))

	fmt.Println("New application successfully created!")
}
