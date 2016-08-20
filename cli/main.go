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
			Usage: raccoon.TASKS_USAGE,
			Action: func(c *cli.Context) error {
				jobParser := parser.JobParser{}
				err := jobParser.CreateJobWithFilePaths(c.String(raccoon.TASKS_NAME),
					c.String(raccoon.INFRASTRUCTURE))
				if err != nil {
					return err
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  raccoon.TASKS_FLAG_ALIAS,
					Usage: raccoon.TASKS_USAGE,
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
