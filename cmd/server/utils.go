package main

import (
	"fmt"
	"os"
)

func dieOnError(err error, msg string) {
	if err == nil {
		return
	}
	fmt.Printf("[ERROR]: %s: %v\n", msg, err)
	os.Exit(1)
}

func getEnvOrDie(name string) string {
	e := os.Getenv(name)
	if e != "" {
		return e
	}
	dieOnError(fmt.Errorf("Variable %s was not set", name), "Missing required environment varialbe")
	return ""
}
