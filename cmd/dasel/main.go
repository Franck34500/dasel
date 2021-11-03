package main

import (
	"fmt"
	"github.com/tomwright/dasel/internal/command"
	"os"
)

func main() {
	cmd := command.NewRootCMD()
	command.ChangeDefaultCommand(cmd, "select", "-v", "--version", "help")
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error: "+err.Error())
		os.Exit(1)
	}
}
