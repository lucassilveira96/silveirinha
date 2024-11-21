package main

import (
	"log"
	"os"
	"silveirinha/cmd"
)

func main() {
	// Calls the Execute function to run the root command
	if err := cmd.Execute(); err != nil {
		// If an error occurs, display the message and exit the program with a non-zero code
		log.Println("Error:", err)
		os.Exit(1)
	}
}
