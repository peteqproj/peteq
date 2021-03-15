package main

import (
	_ "github.com/lib/pq"
	cmd "github.com/peteqproj/peteq/cmd/server/cmd"
)

func main() {
	cmd.Execute()
}
