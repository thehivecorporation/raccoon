package constants

/* GLOBAL */

const (
	VERSION         string = "0.2.1"
	APP_DESCRIPTION string = "WIP App orchestration, configuration and deployment"
	APP_NAME        string = "Raccoon"
)

/* CLI COMMANDS */

//INSTRUCTIONS AND HOSTS
const (
	INSTRUCTIONS_NAME string = "zombiebook"

	INSTRUCTIONS_FLAG_NAME  string = "zombiebook"
	INSTRUCTIONS_FLAG_ALIAS string = "zombiebook, z"
	INSTRUCTIONS_USAGE      string = "Execute a Zombiebook"

	HOSTS_FLAG_NAME     string = "mansion"
	HOSTS_NAME          string = "room"
	HOST_LAUNCH_MESSAGE string = "Entering " + HOSTS_NAME + "..."
	HOSTS_FLAG_ALIAS    string = "mansion, m"
	HOSTS_FLAG_USAGE    string = "The Mansion file"

	RELATIONSHIP_KEY string = "chapter"

	GROUP_NAME string = "mansion"
	MAINTAINER string = "maintainer"
)

//SERVER
const (
	SERVER_NAME  = "server"
	SERVER_USAGE = "Launch a server to receive Zombiebook JSON files"

	PORT_FLAG_NAME  = "port"
	PORT_FLAG_ALIAS = "port, p"
	PORT_FLAG_USAGE = "The port to run the server on"
)
