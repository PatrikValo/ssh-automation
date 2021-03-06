package main

import (
	"os"

	"github.com/PatrikValo/ssh-automation/pkg/cli"
	"github.com/PatrikValo/ssh-automation/pkg/core"
	"github.com/PatrikValo/ssh-automation/pkg/printer"
	"github.com/PatrikValo/ssh-automation/pkg/program"
)

func main() {
	commandLine := cli.CreateCli()
	userInput, err := commandLine.GetUserInput()

	if err != nil {
		panic(err)
	}

	parser := program.Parser{Filename: userInput.Filename}

	parsedProgram, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	user := userInput.Auth.User
	auth, err := userInput.Auth.GetAuthMethod()
	if err != nil {
		panic(err)
	}

	terminal := make(chan printer.Printer)
	executor := core.ExecutorFactory(&parsedProgram.Config, parsedProgram.Hosts, user, &auth, terminal)

	go func() {
		defer close(terminal)
		err := executor.Execute(parsedProgram)
		if err != nil {
			os.Exit(1)
		}
	}()

	printer.RegisterPrinter(terminal)
	os.Exit(0)
}
