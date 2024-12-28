package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Get environment from arguments or use "development" as default
	env := "development"
	if len(os.Args) > 1 {
		env = os.Args[1]
	}

	config, err := LoadConfig(env)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	fmt.Printf("Loaded Configuration for %s:\n", env)
	fmt.Printf("App Name: %s\n", config.AppName)
	fmt.Printf("Port: %d\n", config.Port)
	fmt.Printf("Debug Mode: %v\n", config.Debug)
}
