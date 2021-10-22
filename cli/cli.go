package cli

import (
	"fmt"
	"github.com/techcraftlabs/base"
	"github.com/techcraftlabs/base/io"
	clix "github.com/urfave/cli/v2"
	"os"
)

type (
	App struct {
		app *clix.App
		baseHttp  *base.Client
	}
)

func New() *App {

	app := new(App)

	desc :=
		`regctl is a command line tool intended for managing access to igridnet it manages users, nodes and regions`

	author1 := &clix.Author{
		Name:  "Pius Alfred",
		Email: "me.pius1102@gmail.com",
	}

	 cliApp := new(clix.App)

	cliApp = &clix.App{
		Name:  "regctl",
		Usage: "commandline tool for igridnet identity and access management",
		Version:              "1.0.0",
		Description:          desc,
		Commands:             []*clix.Command{

		},
		Flags:                flags(),
		EnableBashCompletion: true,
		Before:               beforeActionFunc,
		After:                afterActionFunc,
		CommandNotFound:      onCommand404,
		OnUsageError:         onErrFunc,
		Authors:              authors(author1),
		Copyright:            "MIT Licence, Creative Commons",
		ErrWriter:            os.Stderr,
	}


	app.baseHttp = base.NewClient()

	cliApp.Commands = []*clix.Command{
		app.adminsCommand(),
		app.nodesCommand(),
		app.regionsCommand(),
	}
	app.app = cliApp

	return app
}

func (app *App)Run(args []string)error{
	return app.app.Run(args)
}

func beforeActionFunc(context *clix.Context) error {
	return nil
}

func afterActionFunc(context *clix.Context) error {
	return nil
}

func onCommand404(context *clix.Context, s string) {
	_, _ = fmt.Fprintf(io.Stderr, "not found: %s\n", s)
}

func onErrFunc(context *clix.Context, err error, subcommand bool) error {
	_, _ = fmt.Fprintf(io.Stderr, "error: %v\n", err)
	return nil
}



func appendCommands(comm ...*clix.Command) []*clix.Command {
	var commands []*clix.Command
	for _, command := range comm {
		commands = append(commands, command)
	}
	return commands
}

func flags(fs ...clix.Flag) []clix.Flag {
	var flgs []clix.Flag
	for _, flg := range fs {
		flgs = append(flgs, flg)
	}
	return flgs
}

func authors(auth ...*clix.Author) []*clix.Author {
	var authors []*clix.Author
	for _, author := range auth {
		authors = append(authors, author)
	}
	return authors
}

