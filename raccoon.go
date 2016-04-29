package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/dispatcher"
	"github.com/thehivecorporation/raccoon/instructions"
	"github.com/thehivecorporation/raccoon/job"
	"encoding/json"
	"io/ioutil"
	"os"
)

type BOOK struct {
	RUN []instructions.RUN
	ADD []instructions.ADD
}


func main() {
	var file string
	var book BOOK
	node := connection.Node{
		Username: "vagrant",
		Password: "vagrant",
		IP:       "192.168.33.10",
	}

	raccoon_app := cli.NewApp()
	raccoon_app.Name = "Racoon"
	raccoon_app.Usage = "WIP App orchestration, configuration and deployment"
	raccoon_app.Action = func(c *cli.Context) {
		println("\n" +
			"######      #      #####    #####   #######  #######  #     # \n" +
			"#     #    # #    #     #  #     #  #     #  #     #  ##    # \n" +
			"#     #   #   #   #        #        #     #  #     #  # #   # \n" +
			"######   #     #  #        #        #     #  #     #  #  #  # \n" +
			"#   #    #######  #        #        #     #  #     #  #   # # \n" +
			"#    #   #     #  #     #  #     #  #     #  #     #  #    ## \n" +
			"#     #  #     #   #####    #####   #######  #######  #     # \n")
	}

	raccoon_app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "file",
			Value:       "JSON",
			Usage:       "Read JSON Zombiebook",
			Destination: &file,
		},
	}


	raccoon_app.Action = func(c *cli.Context) {
		log.WithFields(log.Fields{
			"File": file,
		}).Info("READING---------------------------------------> ")

	}
	raccoon_app.Run(os.Args)


	dat, err := ioutil.ReadFile(file)
	if err!=nil{
		log.Error(err)
	}
	log.Info(string(dat))
	json.Unmarshal([]byte(string(dat)), &book)
	log.Info(book)
	nodes := make([]connection.Node, 1)
	nodes[0] = node

	cluster := connection.Cluster{
		Nodes: nodes,
	}

	demo_instructions := make([]instructions.Instruction, 2)
	demo_instructions[0] = &instructions.RUN{"RUN", "Install EPEL repo", "sudo yum install -y epel"}
	demo_instructions[1] = &instructions.RUN{"RUN", "Install tar", "sudo yum install -y tar"}

	recipe := job.Recipe{
		Instructions: demo_instructions,
	}

	job := job.Job{
		Cluster: cluster,
		Recipe:  recipe,
	}

	dispatcher.Dispatch(job)
}
