package utils

import (
	"fmt"
	"os"
)

// DieOnError exit process when error found
// should be used only from cmd
func DieOnError(err error, msg string) {
	if err == nil {
		return
	}
	fmt.Printf("[ERROR]: %s: %v\n", msg, err)
	os.Exit(1)
}

// GetEnvOrDie returns environment varialbe
// when variables is not set or its empty the process stops
func GetEnvOrDie(name string) string {
	e := os.Getenv(name)
	if e != "" {
		return e
	}
	DieOnError(fmt.Errorf("Variable %s was not set", name), "Missing required environment varialbe")
	return ""
}
