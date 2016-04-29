package constants

/* GLOBAL */

const (
	VERSION         string = "0.0.1"
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
	HOST_NAME           string = "room"
	HOST_LAUNCH_MESSAGE string = "Entering " + HOST_NAME + "..."
	HOSTS_FLAG_ALIAS    string = "mansion, m"
	HOSTS_FLAG_USAGE    string = "The Mansion file"
)

//SERVER
const (
	SERVER_NAME  = "server"
	SERVER_USAGE = "Launch a server to receive Zombiebook JSON files"
)
