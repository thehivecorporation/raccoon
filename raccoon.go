package main

import (
	"./book"
	"github.com/codegangsta/cli"
	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/dispatcher"
	"github.com/thehivecorporation/raccoon/instructions"
	"github.com/thehivecorporation/raccoon/job"
	"os"
)

func main() {
	var file string
	var zbook book.BOOK

	node := connection.Node{
		Username: "vagrant",
		Password: "vagrant",
		IP:       "192.168.33.10",
	}

	raccoon_app := cli.NewApp()
	raccoon_app.Name = "Racoon"
	raccoon_app.Usage = "WIP App orchestration, configuration and deployment"
	raccoon_app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "file, -f",
			Usage:       "Read JSON Zombiebook",
			Destination: &file,
		},
	}
	raccoon_app.Action = func(c *cli.Context) {
		println("\n" +
			"######      #      #####    #####   #######  #######  #     # \n" +
			"#     #    # #    #     #  #     #  #     #  #     #  ##    # \n" +
			"#     #   #   #   #        #        #     #  #     #  # #   # \n" +
			"######   #     #  #        #        #     #  #     #  #  #  # \n" +
			"#   #    #######  #        #        #     #  #     #  #   # # \n" +
			"#    #   #     #  #     #  #     #  #     #  #     #  #    ## \n" +
			"#     #  #     #   #####    #####   #######  #######  #     # \n")

		if len(c.String("file"))!=0{
			zombie_book := &book.BOOK{}
			zbook = zombie_book.ReadZbook(file)
		}

	}

	raccoon_app.Run(os.Args)

	nodes := make([]connection.Node, 1)
	nodes[0] = node

	cluster := connection.Cluster{
		Nodes: nodes,
	}

	demo_instructions := make([]instructions.Instruction, len(zbook.RUN))
	for index,command := range zbook.RUN {
		demo_instructions[index] = &instructions.RUN{command.Name,command.Description,command.Instruction}
	}

	recipe := job.Recipe{
		Instructions: demo_instructions,
	}

	job := job.Job{
		Cluster: cluster,
		Recipe:  recipe,
	}

	dispatcher.Dispatch(job)
}
