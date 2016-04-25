package main

import "log"

type Command struct {
	Name    string
	Command string
}

type Host struct {
	IP string
}

type Node struct {
	IP           string
	Username     string
	Password     string
	AuthFilePath string
}

type Cluster struct {
	Nodes []Node
}

type Recipe struct {
	Commands []Command
}

type Job struct {
	Cluster
	Recipe
}

func main() {
	node := Node{
		Username: "vagrant",
		Password: "vagrant",
		IP:       "192.168.33.10",
	}

	nodes := make([]Node, 1)
	nodes[0] = node

	cluster := Cluster {
		Nodes: nodes,
	}

	command1 := Command{
		Name:    "Install EPEL repo",
		Command: "sudo yum install -y epel",
	}

	command2 := Command{
		Name:    "Install tar",
		Command: "sudo yum install -y tar",
	}

	commands := make([]Command, 2)
	commands[0] = command1
	commands[1] = command2

	recipe := Recipe{
		Commands: commands,
	}

	job := Job{
		Cluster: cluster,
		Recipe: recipe,
	}

	quit := make(chan bool)

	Dispatch(job, quit)

	for _ = range quit {
		log.Println("Finished")
		break
	}
}
