package parser

import "testing"

func TestReadTaskFile(t *testing.T) {
	fileParser := File{}

	//Non existing filepath
	filePath := "/tmp/i-do-not-exist"

	_, err := fileParser.Parse(filePath)
	if err == nil {
		t.Fatal("After passing a non existing file a descriptive error must be thrown")
	}
}
