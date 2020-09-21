package local

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

type (
	// DB is local db
	DB struct {
		Path  string
		mutex sync.Mutex
	}
)

func (d *DB) Read() ([]byte, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return ioutil.ReadFile(d.Path)
}

func (d *DB) Write(data []byte) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if err := ioutil.WriteFile(d.Path, data, os.ModePerm); err != nil {
		return fmt.Errorf("Failed to update local db: %w", err)
	}
	return nil
}
