package main

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
	Recipe Recipe
	Output chan string
}

func main() {
	node := Node{
		Username: "vagrant",
		Password: "vagrant",
		IP:       "192.168.33.10",
	}

	command := Command{
		Name:    "Install EPEL repo",
		Command: "sudo yum install -y epel",
	}

	ExecuteCommandOnNode(command, node)
}
