package parser

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestReadZbookFile(t *testing.T) {
	//Non existing filepath
	filePath := "/tmp/i-do-not-exist"

	_, err := readCommandFile(filePath)
	if err == nil {
		t.Fatal("After passing a non existing file a descriptive error must be thrown")
	}

	//Show error if you use "go test -v"
	t.Log(err)


	//Create a file with a wrong JSON structure / syntax and pass it
	f, _ := ioutil.TempFile("/tmp", "wrongJson")
	f.WriteString("{wrong:\"syntax\",}")
	f.Close()
	filePath = f.Name()

	_, err = readCommandFile(filePath)
	if err == nil {
		t.Fatal("After passing a json with wrong structure an error must be thrown")
	}

	//Show error if you use "go test -v"
	t.Log(err)


	//Pass a correct JSON
	filePath = "../examples/exampleCommands.json"

	_, err = readCommandFile(filePath)
	if err != nil {
		t.Fatalf("JSON file was correct. No error must be here: %s\n", err.Error())
	}

	//Pass a syntactically correct zbook with no content
	f, _ = ioutil.TempFile("/tmp", "no-content")
	c := []Commands{
		Commands{
			Title:           "chapter1",
			Maintainer:      "maintainer",
			Command: make([]map[string]string, 0), //Empty
		},
	}
	by, _ := json.Marshal(c)
	f.Write(by)
	f.Close()
	filePath = f.Name()

	_, err = readCommandFile(filePath)
	if err != nil {
		t.Fatal("No instructions were provided but syntax was correct:", err)
	}

	//Create mocked data
	cs := []Commands{
		Commands{
			Title:      "chapter1",
			Maintainer: "maintainer",
			Command: []map[string]string{
				0: {
					"name":       "RUN",
					"description": "Install htop",
					"instruction": "sudo yum install -y htop",
				},
				1:{
					"name":"ADD",
					"description":"a description",
					"sourcePath":"source path",
					"destPath":"destination path",
				},
			},
		},
	}

	//Clone the data 5 times
	batteryTests := make([][]Commands,5)
	for i := 0; i<5; i++ {
		batteryTests[i] = cs
	}

	//Delete some correct keys to replace them for incorrect ones
	delete(batteryTests[0][0].Command[1], "instruction")
	batteryTests[0][0].Command[1]["instru2ction"] = "sudo yum install -y htop"

	delete(batteryTests[1][0].Command[1], "sourcePath")
	batteryTests[1][0].Command[1]["sourceP2ath"] = "source"

	delete(batteryTests[2][0].Command[1], "destPath")
	batteryTests[2][0].Command[1]["dest2Path"] = "dest"

	delete(batteryTests[3][0].Command[1], "name")
	batteryTests[3][0].Command[1]["n2me"] = "ADD"

	for i := 0; i<5; i++ {
		f, _ = ioutil.TempFile("/tmp", "raccoon_temp")
		by, _ := json.Marshal(batteryTests[i])
		f.Write(by)
		f.Close()
		filePath = f.Name()

		_, err = readCommandFile(filePath)
		if err == nil {
			t.Fatalf("Syntax was incorrect but no error found on index %d",i)
		}
	}
}
