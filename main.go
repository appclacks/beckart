package main

import (
	"github.com/appclacks/beckart/cmd"
)

func main() {
	command := cmd.RunCommand()
	_ = command.Execute()
}
