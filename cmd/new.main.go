package cmd

import "fmt"

func (n S1GU) createMainPage(appName string) string {
	return fmt.Sprintf(`
	package main 
	import ( 
		"%s/cmd"
	) 
	
	func main() { 
		cmd.Execute()
	}`, appName)
}
