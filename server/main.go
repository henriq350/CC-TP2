package main

import (
	"ccproj/utils"
	"fmt"
	"os"
)

func main() {

	// Check if the user provided a configuration file
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <config.json>")
		return
	}

	configFile := os.Args[1]

	if !gUtils.IsJSONFile(configFile) {
		fmt.Printf("Error: Configuration file must be a .json\n")
		return
	}

	// Parse the configuration file
	task, err := ParseTasks(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	PrintTasks(task)
}
