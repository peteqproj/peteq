package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

// JSONStringToReadCloser converts json string into io.ReadCloser
// should be used only in test as this method will exit on error
func JSONStringToReadCloser(j map[string]interface{}) io.ReadCloser {
	b, err := json.Marshal(j)
	DieOnError(err, "Failed to convert json to io.ReadCloser")
	return ioutil.NopCloser(bytes.NewReader(b))
}

// MustMarshal marshals or dies
func MustMarshal(v interface{}) []byte {
	r, err := json.Marshal(v)
	if err != nil {
		DieOnError(err, "Failed to marshal")
	}
	return r
}
