//Package main contains the CLI application of Raccoon along with the options
//of the CLI interface
package main

import (
	"fmt"
	"os"

	"github.com/thehivecorporation/raccoon"
)

// CLI
//
// Raccon
// Copyright (C) 2016 The Hive Corporation
//
// This program comes with ABSOLUTELY NO WARRANTY; for details type 'warranty'.
// This is free software, and you are welcome to redistribute it under certain conditions; read License.md file for details.
//
// NAME:
// Raccoon - WIP Automation utility made easy with Dockerfile syntax
//
// USAGE:
// cli [global options] command [command options] [arguments...]
//
// VERSION:
// 0.3.0
//
// COMMANDS:
// job     Execute a job
// server  Launch a server to receive Commands JSON files
// show    Show special information about Raccoon
//
// GLOBAL OPTIONS:
// --help, -h     show help
// --version, -v  print the version

//
// Job option
//
// NAME:
// cli job - Execute a job
//
// USAGE:
// cli job [command options] [arguments...]
//
// OPTIONS:
// --tasks value, -t value           Tasks file
// --infrastructure value, -i value  Infrastructure file
// --job value, -j value             Job file containing infrastructure and tasks information
// --dispatcher value, -d value      Dispatching method between 3 options: sequential (no concurrent dispatch). simple (a Goroutine for each host) and worker_pool (a fixed number of workers) (default: "simple")
// --workersNumber value, -w value   In case of worker_pool dispath method, define the maximum number of workers here (default: 5)
//
// For example:
//	raccoon tasks -t tasksFile.json -i infrastructureFile.json
//
// Server command
//
//	NAME:
//	raccoon server - Launch a server to receive Commands JSON files
//
//	USAGE:
//	raccoon server [command options] [arguments...]
//
//	OPTIONS:
//	--port value  port, p (default: "8123")
//
// For example:
//	raccoon server -p 8080

import (
	"strings"

	"github.com/codegangsta/cli"
	"github.com/thehivecorporation/raccoon/parser"
	"github.com/thehivecorporation/raccoon/server"
)

func dispatcherFactory(option string, workersNumber int) raccoon.Dispatcher {
	switch option {
	default:
		return new(raccoon.SimpleDispatcher)
	case "sequential":
		return new(raccoon.SequentialDispatcher)
	case "workers_pool":
		return &raccoon.WorkerPoolDispatcher{
			Workers: workersNumber,
		}
	}
}

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
			Name:  "job",
			Usage: "Execute a job",
			Action: func(c *cli.Context) error {
				jobParser := parser.Job{}
				jobParser.Dispatcher = dispatcherFactory(c.String("dispatcher"), c.Int("workersNumber"))

				//Parse a full job file if applicable
				if c.String("job") != "" {
					genericParser := parser.Generic{}
					jobFile, err := genericParser.Parse(c.String("job"))
					if err != nil {
						return err
					}

					var req raccoon.JobRequest
					if err := genericParser.Build(jobFile, &req); err != nil {
						return err
					}

					taskList, err := jobParser.ParseTaskList(req.TaskList)
					if err != nil {
						return err
					}

					infParser := parser.InfrastructureFile{}
					infParser.Prepare(req.Infrastructure)

					jobs := jobParser.BuildJobList(req.Infrastructure, taskList)

					jobParser.Dispatcher.Dispatch(*jobs)

					return nil
				}

				if c.String("single-task") != "" {
					infFilePath := c.String("infrastructure")
					taskFilePath := c.String("tasks")

					jobParser := parser.Job{}
					jobParser.Dispatcher = dispatcherFactory(c.String("dispatcher"), c.Int("workersNumber"))

					genericParser := parser.Generic{}

					var infrastructure raccoon.Infrastructure
					if err := genericParser.ParserFactory(infFilePath, &infrastructure); err != nil {
						return err
					}

					taskList := new([]raccoon.Task)
					if err := genericParser.ParserFactory(taskFilePath, taskList); err != nil {
						return err
					}
					taskList, err := jobParser.ParseTaskList(taskList)
					if err != nil {
						return err
					}

					params := strings.Split(c.String("single-task"), "=")
					if len(params) != 2 {
						return fmt.Errorf("When specifying a single task you " +
							"must provide a cluster name and a task name " +
							"separated by a '='\nFor example: --single-task " +
							"cluster1=task1")
					}

					relationList := raccoon.RelationList{
						{
							ClusterName: params[0],
							TaskList: []string{
								params[1],
							},
						},
					}

					relationParser := parser.Relation{}
					if err := relationParser.Prepare(&infrastructure, &relationList); err != nil {
						return err
					}

					infParser := parser.InfrastructureFile{}
					infParser.Prepare(&infrastructure)

					jobs := jobParser.BuildJobList(&infrastructure, taskList)

					jobParser.Dispatcher.Dispatch(*jobs)

					return nil
				}

				//Parse the 3 relation files
				if c.String("relation") != "" {
					infFilePath := c.String("infrastructure")
					taskFilePath := c.String("tasks")
					relationFilePath := c.String("relation")

					jobParser := parser.Job{}
					jobParser.Dispatcher = dispatcherFactory(c.String("dispatcher"), c.Int("workersNumber"))

					genericParser := parser.Generic{}

					var infrastructure raccoon.Infrastructure
					if err := genericParser.ParserFactory(infFilePath, &infrastructure); err != nil {
						return err
					}

					taskList := new([]raccoon.Task)
					if err := genericParser.ParserFactory(taskFilePath, taskList); err != nil {
						return err
					}
					taskList, err := jobParser.ParseTaskList(taskList)
					if err != nil {
						return err
					}

					relationParser := parser.Relation{}
					var relationList raccoon.RelationList
					if err := relationParser.ParserFactory(relationFilePath, &relationList); err != nil {
						return err
					}
					if err := relationParser.Prepare(&infrastructure, &relationList); err != nil {
						return err
					}

					infParser := parser.InfrastructureFile{}
					infParser.Prepare(&infrastructure)

					jobs := jobParser.BuildJobList(&infrastructure, taskList)

					jobParser.Dispatcher.Dispatch(*jobs)

					return nil

				}

				if err := jobParser.CreateJobWithFilePaths(c.String("tasks"),
					c.String("infrastructure")); err != nil {
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
					Name:  "infrastructure, i",
					Usage: "Infrastructure file",
				},
				cli.StringFlag{
					Name:  "job, j",
					Usage: "Job file containing infrastructure and tasks information",
				},
				cli.StringFlag{
					Name:  "relation, r",
					Usage: "A relation file that points to an infrastructure and tasks files",
				},
				cli.StringFlag{
					Name:  "single-task, s",
					Usage: "[clusterName]=[taskName] Use single-task when you " +
						"just want to execute a unique task in a unique cluster. " +
						"Write the name of the cluster referenced on the " +
						"infrastructure file on the left side of the equal and " +
						"the task referenced on the task list name on the right side",
				},
				cli.StringFlag{
					Name: "dispatcher, d",
					Usage: "Dispatching method between 3 options: sequential " +
						"(no concurrent dispatch). simple (a Goroutine for " +
						"each host) and worker_pool (a fixed number of workers)",
					Value: "simple",
				},
				cli.IntFlag{
					Name:  "workersNumber, w",
					Usage: "In case of worker_pool dispath method, define the maximum number of workers here",
					Value: 5,
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
			Name:  "show",
			Usage: "Show special information about Raccoon",
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
