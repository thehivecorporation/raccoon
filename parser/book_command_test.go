package parser

import "testing"

func TestExecuteZombieBook(t *testing.T) {
	//No zbook specified
	err := ExecuteZombieBook("", "examples/exampleMansion.json")
	if err == nil {
		t.Fatal("An error must be thrown when no filepath is specified in a zombiebook")
	}

	//No mansion specified
	err = ExecuteZombieBook("examples/exampleBook.json", "")
	if err == nil {
		t.Fatal("An error must be thrown when no filepath is specified in a mansion")
	}

	//Zombiebook file doesn't exists
	err = ExecuteZombieBook("/tmp/i-do-not-exist", "examples/exampleMansion.json")
	if err == nil {
		t.Fatal("An error must be thrown when zombiebook file doesn't exist")
	}

	//Mansion file doesn't exists
	err = ExecuteZombieBook("../examples/exampleBook.json", "/tmp/i-do-not-exist")
	if err == nil {
		t.Fatal("An error must be thrown when mansion file doesn't exist")
	}
}
