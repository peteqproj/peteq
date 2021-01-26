package cmd

import "github.com/peteqproj/peteq/pkg/utils"

// DieOnError kills the process and prints a message
func DieOnError(err error, msg string) {
	utils.DieOnError(err, msg)
}
