package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/thehivecorporation/raccoon/parser"
	"github.com/thehivecorporation/raccoon/server"
	"github.com/thehivecorporation/raccoon/constants"
)

//Init initializes the CLI functions
func main() {
	app := cli.NewApp()
	app.Name = constants.APP_NAME
	app.Usage = constants.APP_DESCRIPTION
	app.Version = constants.VERSION

	app.Commands = []cli.Command{
		{
			Name:  constants.COMMANDS_NAME,
			Usage: constants.COMMANDS_USAGE,
			Action: func(c *cli.Context) error {
				err := parser.ExecuteCommandsFile(c.String(constants.COMMANDS_NAME),
					c.String(constants.HOSTS_FLAG_NAME))
				if err != nil {
					return err
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  constants.COMMANDS_FLAG_ALIAS,
					Usage: constants.COMMANDS_USAGE,
				},
				cli.StringFlag{
					Name: constants.HOSTS_FLAG_ALIAS,

					Usage: constants.HOSTS_FLAG_USAGE,
				},
			},
		},
		{
			Name:   constants.SERVER_NAME,
			Usage:  constants.SERVER_USAGE,
			Action: server.REST,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  constants.PORT_FLAG_ALIAS,
					Usage: constants.PORT_FLAG_USAGE,
					Value: "8123",
				},
			},
		},
	}

	app.Run(os.Args)
}
