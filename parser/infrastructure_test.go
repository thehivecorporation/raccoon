package parser

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/thehivecorporation/raccoon"
)

func CreateTestFileWithString(content string) (filePath string) {
	f, _ := ioutil.TempFile("/tmp", "raccoon")
	f.WriteString(content)
	f.Close()
	filePath = f.Name()

	return
}

func TestInfrastructureFile(t *testing.T) {
	content := []byte(`{"name": "A name","infrastructure": [ {"name": "some cluster","tasks": ["task2"],"hosts": [{ "ip": "172.17.42.1", "sshPort": 32768, "username": "root", "description": "cassandra01", "password": "root"},{ "ip": "172.17.42.1", "sshPort": 32769, "description": "cassandra02", "username": "root", "interactiveAuth": true},{ "ip": "172.17.42.1", "sshPort": 32768, "description": "cassandra03", "username": "root", "password": "root"}] }]},"tasks": [{ "title": "task1", "maintainer": "Burkraith", "commands": [{"name": "ADD","sourcePath": "doc.go","destPath": "/tmp","description": "Raccoon.go to /tmp"},{"name": "RUN","description": "Removing htop","instruction": "sudo yum remove -y htop"},{"name": "ADD","sourcePath": "main.go","destPath": "/tmp","description": "copying raccoon.go to /tmp"} ]},{ "title": "task2", "maintainer": "Mario", "commands": [{"name": "RUN","description": "Removing htop","instruction": "sudo apt-get remove -y htop"} ]}] }`)
	r := bytes.NewReader(content)
	infParser := InfrastructureFile{}
	var inf raccoon.Infrastructure
	if err := infParser.Build(r, &inf); err != nil {
		t.Fatal(err)
	}

	if inf.Name != "A name" {
		t.Error("Expected name 'A name' wasn't found")
	}
}
