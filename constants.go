package raccoon

/* GLOBAL */

const (
	VERSION         string = "0.2.2"
	APP_DESCRIPTION string = "WIP App orchestration, configuration and deployment"
	APP_NAME        string = "Raccoon"
)

/* CLI COMMANDS */

//INSTRUCTIONS AND HOSTS
const (
	TASKS_NAME string = "tasks"

	TASKS_FLAG_NAME  string = "tasks"
	TASKS_FLAG_ALIAS string = "tasks, t"
	TASKS_USAGE      string = "Execute a task list"

	INFRASTRUCTURE            string = "infrastructure"
	HOSTS_NAME                string = "cluster"
	HOST_LAUNCH_MESSAGE       string = "Entering " + HOSTS_NAME + "..."
	INFRASTRUCTURE_FLAG_ALIAS string = "infrastructure, i"
	INFRASTRUCTURE_FLAG_USAGE string = "The Infrastructure file"

	RELATIONSHIP_KEY string = "commands"

	GROUP_NAME string = "infrastructure"
	MAINTAINER string = "maintainer"
)

//SERVER
const (
	SERVER_NAME  = "server"
	SERVER_USAGE = "Launch a server to receive Commands JSON files"

	PORT_FLAG_NAME  = "port"
	PORT_FLAG_ALIAS = "port, p"
	PORT_FLAG_USAGE = "The port to run the server on"
)
