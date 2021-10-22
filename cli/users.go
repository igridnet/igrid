package cli

import "github.com/urfave/cli/v2"

func (app *App) adminsCommand() *cli.Command {
	return &cli.Command{
		Name:        "admins",
		Usage:       "admins management",
		Description: "admins management: this command let you perform CRUD operation on the registry about admins",
		Subcommands: []*cli.Command{
			app.adminsAddCommand(),
			app.adminLoginCommand(),
		},
	}
}

func (app *App) adminsAddCommand() *cli.Command {
	return &cli.Command{
		Name:        "add",
		Usage:       "add new admin",
		Description: "add: this command let you add new admin in the registry",
		Subcommands: nil,
	}
}

func (app *App) adminLoginCommand() *cli.Command {
	return &cli.Command{
		Name:        "login",
		Usage:       "get access token",
		Description: "enter username and password to receive admin access key",
		Subcommands: nil,
	}
}
