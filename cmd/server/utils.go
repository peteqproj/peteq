package main

import (
	"fmt"
	"os"
)

func dieOnError(err error, msg string) {
	if err == nil {
		return
	}
	fmt.Printf("[ERROR]: %s: %w", msg, err)
	os.Exit(1)
}
