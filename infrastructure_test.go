package raccoon

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/thehivecorporation/raccoon/connection"
	"github.com/thehivecorporation/raccoon"
)

//TODO fix this test
func CreateTestFileWithJSON(content interface{}) (filePath string) {
	by, _ := json.Marshal(content)

	f, _ := ioutil.TempFile("/tmp", "raccoon")
	f.Write(by)
	f.Close()
	filePath = f.Name()

	return
}

func CreateTestFileWithString(content string) (filePath string) {
	f, _ := ioutil.TempFile("/tmp", "raccoon")
	f.WriteString(content)
	f.Close()
	filePath = f.Name()

	return
}

func TestReadMansionFile(t *testing.T) {
	//Create syntactically correct mock object without instructions
	rawContent := Infrastructure{
		Name: "mansion name",
		Infrastructure: []Cluster{
			Cluster{
				Name:     "room name",
				Task: "a chapter",
				Hosts: []Host{
					Host{
						IP:       "192.168.1.44",
						Username: "vagrant",
						Password: "vagrant",
					},
					Host{
						IP:       "192.168.1.45",
						Username: "vagrant",
						Password: "vagrant",
					},
				},
			},
		},
	}

	//Check mansion file doesn't exist
	filePath := "/tmp/i-do-not-exist"
	_, err := readInfrastructureFile(filePath)
	if err == nil {
		t.Fatal("An error must be thrown when passing a non existing file")
	}

	//Check mock object
	_, err = readInfrastructureFile(CreateTestFileWithJSON(rawContent))
	if err != nil {
		t.Fatal(err)
	}

	//Pass a wrong syntax to json
	_, err = readInfrastructureFile(CreateTestFileWithString("{wrong:\"syntax\",}"))
	if err == nil {
		t.Fatal("An error must be thrown when using a wrong JSON syntax")
	}

	//Pass syntactically correct data with wrong mansion syntax

	testNumber := 8
	wrongSyntax := make([]string, testNumber)

	//Empty rooms
	wrongSyntax[0] = "{\"name\":\"A name\", \"rooms\":[] }"

	//Empty hosts
	wrongSyntax[1] = "{\"name\":\"A name\", \"rooms\":[{\"name\":\"some room\", \"chapter\":\"chapter1\", \"hosts\":[] } ] }"

	//Incorrect room name key: "n4me'
	wrongSyntax[2] = "{\"name\":\"A name\", \"rooms\":[{\"n4me\":\"some room\", \"chapter\":\"chapter1\", \"hosts\":[{\"ip\":\"192.168.1.44\", \"username\":\"vagrant\", \"password\":\"vagrant\"} ] } ] }"

	//Incorrect room chapter key: "ch4pter"
	wrongSyntax[3] = "{\"name\":\"A name\", \"rooms\":[{\"name\":\"some room\", \"ch4pter\":\"chapter1\", \"hosts\":[{\"ip\":\"192.168.1.44\", \"username\":\"vagrant\", \"password\":\"vagrant\"} ] } ] }"

	//Incorrect ip key on host: "1p"
	wrongSyntax[4] = "{\"name\":\"A name\", \"rooms\":[{\"name\":\"some room\", \"chapter\":\"chapter1\", \"hosts\":[{\"1p\":\"192.168.1.44\", \"username\":\"vagrant\", \"password\":\"vagrant\"} ] } ] }"

	//Incorrect username key on host: "user"
	wrongSyntax[5] = "{\"name\":\"A name\", \"rooms\":[{\"name\":\"some room\", \"chapter\":\"chapter1\", \"hosts\":[{\"ip\":\"192.168.1.44\", \"user\":\"vagrant\", \"password\":\"vagrant\"} ] } ] }"

	//Incorrect password key on host: "pass"
	wrongSyntax[6] = "{\"name\":\"A name\", \"rooms\":[{\"name\":\"some room\", \"chapter\":\"chapter1\", \"hosts\":[{\"ip\":\"192.168.1.44\", \"username\":\"vagrant\", \"pass\":\"vagrant\"} ] } ] }"

	//Incorrect mansion name
	wrongSyntax[7] = "{\"n4me\":\"A name\", \"rooms\":[{\"name\":\"some room\", \"chapter\":\"chapter1\", \"hosts\":[{\"ip\":\"192.168.1.44\", \"username\":\"vagrant\", \"password\":\"vagrant\"} ] } ] }"

	//TODO incorrectly formatted IP

	for i := 0; i < testNumber; i++ {
		_, err = readInfrastructureFile(CreateTestFileWithString(wrongSyntax[i]))
		if err == nil {
			t.Fatalf("An error must be thrown when passing a incorrect syntax like 'inst2ruction' on index %d\n", i)
		} else {
			t.Log(err)
		}
	}
}
