package cmd

import "fmt"

func CreateMainPage(appName string) string {
	return fmt.Sprintf(`
	package main 
	import ( 
		"%s/cmd"
	) 
	
	func main() { 
		cmd.Execute()
	}`, appName)
}
