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
			Name:  "tasks",
			Usage: "Execute a task list",
			Action: func(c *cli.Context) error {
				jobParser := parser.JobParser{}
				err := jobParser.CreateJobWithFilePaths(c.String("tasks"),
					c.String("infrastructure"))
				if err != nil {
					return err
				}
				return nil
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "tasks, t",
					Usage: "Execute a task list",
				},
				cli.StringFlag{
					Name:  "infrastructure, i",
					Usage: "The Infrastructure file",
				},
			},
		},
		{
			Name:   "server",
			Usage:  "Launch a server to receive Commands JSON files",
			Action: server.REST,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Usage: "port, p",
					Value: "8123",
				},
			},
		},
	}

	app.Run(os.Args)
}
