package cli

import "github.com/urfave/cli/v2"

func (app *App) nodesCommand() *cli.Command {
	return &cli.Command{
		Name:        "nodes",
		Usage:       "nodes management",
		Description: "nodes management: this command let you perform CRUD operation on the registry about nodes",
		Subcommands: []*cli.Command{
			app.nodesAddCommand(),
			app.nodesGetCommand(),
			app.nodesListCommand(),
		},
	}
}

func (app *App) nodesAddCommand() *cli.Command {
	return &cli.Command{
		Name:        "add",
		Usage:       "add new node",
		Description: "add node: this command let you add new node in the registry",
		Subcommands: nil,
	}
}

func (app *App) nodesGetCommand() *cli.Command {
	return &cli.Command{
		Name:        "get",
		Usage:       "get a node by id",
		Description: "get node, get the registered node by specifying its id",
		Subcommands: nil,
	}
}

func (app *App) nodesListCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "list all nodes in the registry",
		Description: "list nodes, list all the registered nodes in the registry",
		Subcommands: nil,
	}
}
