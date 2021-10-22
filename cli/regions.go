package cli

import "github.com/urfave/cli/v2"

func (app *App) regionsCommand() *cli.Command {
	return &cli.Command{
		Name:        "regions",
		Usage:       "regions management",
		Description: "regions management: this command let you perform CRUD operation on the registry about regions",
		Subcommands: []*cli.Command{
			app.regionsAddCommand(),
			app.regionsGetCommand(),
			app.regionsListCommand(),
		},
	}
}

func (app *App) regionsAddCommand() *cli.Command {
	return &cli.Command{
		Name:        "add",
		Usage:       "add new region",
		Description: "add region: this command let you add new region in the registry",
		Subcommands: nil,
	}
}

func (app *App) regionsGetCommand() *cli.Command {
	return &cli.Command{
		Name:        "get",
		Usage:       "get a region by id",
		Description: "get region, get the registered region by specifying its id",
		Subcommands: nil,
	}
}

func (app *App) regionsListCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Usage:       "list all regions in the registry",
		Description: "list regions, list all the registered regions in the registry",
		Subcommands: nil,
	}
}
