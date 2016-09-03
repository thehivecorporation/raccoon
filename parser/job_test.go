package parser

import "testing"

func TestExecuteZombieBook(t *testing.T) {
	jobParser := Job{}

	//No zbook specified
	err := jobParser.CreateJobWithFilePaths("", "examples/exampleInfrastructure.json")
	if err == nil {
		t.Fatal("An error must be thrown when no filepath is specified in a zombiebook")
	}

	//No mansion specified
	err = jobParser.CreateJobWithFilePaths("examples/exampleTasks.json", "")
	if err == nil {
		t.Fatal("An error must be thrown when no filepath is specified in a mansion")
	}

	//Zombiebook file doesn't exists
	err = jobParser.CreateJobWithFilePaths("/tmp/i-do-not-exist", "examples/exampleInfrastructure.json")
	if err == nil {
		t.Fatal("An error must be thrown when zombiebook file doesn't exist")
	}

	//Mansion file doesn't exists
	err = jobParser.CreateJobWithFilePaths("../examples/exampleTasks.json", "/tmp/i-do-not-exist")
	if err == nil {
		t.Fatal("An error must be thrown when mansion file doesn't exist")
	}
}
