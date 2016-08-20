package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/thehivecorporation/raccoon"
	"github.com/thehivecorporation/raccoon/parser"
	"github.com/thehivecorporation/raccoon/server"
)

//Init initializes the CLI functions
func main() {
	app := cli.NewApp()
	app.Name = raccoon.APP_NAME
	app.Usage = raccoon.APP_DESCRIPTION
	app.Version = raccoon.VERSION

	app.Commands = []cli.Command{
		{
			Name:  raccoon.TASKS_NAME,
			Usage: raccoon.COMMANDS_USAGE,
			Action: func(c *cli.Context) error {
				err := parser.CreateJobWithFiles(c.String(raccoon.TASKS_NAME),
					c.String(raccoon.INFRASTRUCTURE))
				if err != nil {
					return err
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  raccoon.COMMANDS_FLAG_ALIAS,
					Usage: raccoon.COMMANDS_USAGE,
				},
				cli.StringFlag{
					Name:  raccoon.INFRASTRUCTURE_FLAG_ALIAS,
					Usage: raccoon.INFRASTRUCTURE_FLAG_USAGE,
				},
			},
		},
		{
			Name:   raccoon.SERVER_NAME,
			Usage:  raccoon.SERVER_USAGE,
			Action: server.REST,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  raccoon.PORT_FLAG_ALIAS,
					Usage: raccoon.PORT_FLAG_USAGE,
					Value: "8123",
				},
			},
		},
	}

	app.Run(os.Args)
}
