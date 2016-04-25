package main

import (
	"time"
	"sync"
)

var wg sync.WaitGroup

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
	node2 := Node{
		Username: "vagrant",
		Password: "vagrant",
		IP:       "192.168.33.11",
	}

	nodes := make([]Node, 2)
	nodes[0] = node
	nodes[1] = node2

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

	command3 := Command{
		Name:    "Showing root path files and folders",
		Command: "ls -la /",
	}

	commands := make([]Command, 3)
	commands[0] = command1
	commands[1] = command2
	commands[2] = command3

	recipe := Recipe{
		Commands: commands,
	}

	job := Job{
		Cluster: cluster,
		Recipe: recipe,
	}


	Dispatch(job)

	time.Sleep(15 * time.Second)

	wg.Wait()
}
