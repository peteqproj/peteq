package main

import (
	_ "github.com/lib/pq"
	cmd "github.com/peteqproj/peteq/cmd/peteq-cli/cmd"
)

func main() {
	cmd.Execute()
}
