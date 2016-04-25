package main

import (
	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon/dispatcher"
	"github.com/thehivecorporation/raccoon/instructions"
	"github.com/thehivecorporation/raccoon/job"
)


func main() {
	node := connection.Node{
		Username: "vagrant",
		Password: "vagrant",
		IP:       "192.168.33.10",
	}

	nodes := make([]connection.Node, 1)
	nodes[0] = node

	cluster := connection.Cluster{
		Nodes: nodes,
	}

	demo_instructions := make([]instructions.Instruction, 2)
	demo_instructions[0] = &instructions.CMD{"CMD","Install EPEL repo","sudo yum install -y epel"}
	demo_instructions[1] = &instructions.CMD{"CMD", "Install tar", "sudo yum install -y tar"}

	recipe := job.Recipe{
		Instructions: demo_instructions,
	}

	job := job.Job{
		Cluster: cluster,
		Recipe:  recipe,
	}

	dispatcher.Dispatch(job)
}
