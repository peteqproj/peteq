package main

import (
	_ "github.com/lib/pq"
	cmd "github.com/peteqproj/peteq/cmd/dev-cli/cmd"
)

func main() {
	cmd.Execute()
}
