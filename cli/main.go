//Package main contains the CLI application of Raccoon along with the options
//of the CLI interface
//
// CLI
//
//	NAME:
//	Raccoon - WIP Automation utility made easy with Dockerfile syntax
//
//	USAGE:
//	cli [global options] command [command options] [arguments...]
//
//	VERSION:
//	0.2.2
//
//	COMMANDS:
//	tasks   Execute a task list
//	server  Launch a server to receive Commands JSON files
//	show    Show special information about Raccoon
//
//	GLOBAL OPTIONS:
//	--help, -h     show help
//	--version, -v  print the version
//
// Tasks command
//
//	NAME:
//	cli tasks - Execute a task list
//
//	USAGE:
//	cli tasks [command options] [arguments...]
//
//	OPTIONS:
//	--tasks value, -t value                     Tasks file
//	--infrastructure value, --warranty i value  Infrastructure file
//
// For example:
//	raccoon tasks -t tasksFile.json -i infrastructureFile.json
//
// Server command
//
//	NAME:
//	cli server - Launch a server to receive Commands JSON files
//
//	USAGE:
//	cli server [command options] [arguments...]
//
//	OPTIONS:
//	--port value  port, p (default: "8123")
//
// For example:
//	raccoon server -p 8080
package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/thehivecorporation/raccoon"
	"github.com/thehivecorporation/raccoon/parser"
	"github.com/thehivecorporation/raccoon/server"
	"fmt"
)

	func main() {
	fmt.Printf("\nRaccon\nCopyright (C) 2016 The Hive Corporation\n\nThis program " +
	"comes with ABSOLUTELY NO WARRANTY; for details type 'warranty'.\nThis is " +
	"free software, and you are welcome to redistribute it under certain " +
	"conditions; read License.md file for details.\n\n")

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
					Usage: "Tasks file",
				},
				cli.StringFlag{
					Name:  "infrastructure,warranty i",
					Usage: "Infrastructure file",
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
		{
			Name:   "show",
			Usage:  "Show special information about Raccoon",
			Action: func(c *cli.Context) error {
				if !c.Bool("warranty") {
					fmt.Println(
`THERE IS NO WARRANTY FOR THE PROGRAM, TO THE
EXTENT PERMITTED BY APPLICABLE LAW. EXCEPT WHEN OTHERWISE
STATED IN WRITING THE COPYRIGHT HOLDERS AND/OR OTHER PARTIES
PROVIDE THE PROGRAM "AS IS" WITHOUT WARRANTY OF ANY KIND,
EITHER EXPRESSED OR IMPLIED, INCLUDING, BUT NOT LIMITED TO,
THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A
PARTICULAR PURPOSE. THE ENTIRE RISK AS TO THE QUALITY AND
PERFORMANCE OF THE PROGRAM IS WITH YOU. SHOULD THE PROGRAM
PROVE DEFECTIVE, YOU ASSUME THE COST OF ALL NECESSARY
SERVICING, REPAIR OR CORRECTION.`)
				}

				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "warranty",
					Usage: "warranty, w",
				},
				cli.BoolFlag{
					Name:  "conditions",
					Usage: "conditions, c",
				},
			},
		},
	}

	app.Run(os.Args)
}
