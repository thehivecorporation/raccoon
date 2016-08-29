package server

import (
	"bytes"
	"testing"
)

func TestParseRequest(t *testing.T) {
	content := []byte(`{"infrastructure": {"name": "A name","infrastructure": [ {"name": "some cluster","tasks": ["task2"],"hosts": [{ "ip": "172.17.42.1", "sshPort": 32768, "username": "root", "description": "cassandra01", "password": "root"},{ "ip": "172.17.42.1", "sshPort": 32769, "description": "cassandra02", "username": "root", "interactiveAuth": true},{ "ip": "172.17.42.1", "sshPort": 32768, "description": "cassandra03", "username": "root", "password": "root"}] }]},"tasks": [{ "title": "task1", "maintainer": "Burkraith", "commands": [{"name": "ADD","sourcePath": "doc.go","destPath": "/tmp","description": "Raccoon.go to /tmp"},{"name": "RUN","description": "Removing htop","instruction": "sudo yum remove -y htop"},{"name": "ADD","sourcePath": "main.go","destPath": "/tmp","description": "copying raccoon.go to /tmp"} ]},{ "title": "task2", "maintainer": "Mario", "commands": [{"name": "RUN","description": "Removing htop","instruction": "sudo apt-get remove -y htop"} ]}] }`)
	r := bytes.NewReader(content)
	_, err := parseRequest(r)
	if err != nil {
		t.Fatal(err)
	}

	content = []byte("{wrong:\"syntax\",}")
	r = bytes.NewReader(content)
	_, err = parseRequest(r)
	if err == nil {
		t.Fatal("An error must be thrown when using a wrong JSON syntax")
	}

	content = []byte(`{"infrastructure": {"name": "A name","infras2tructure": [ {"name": "some cluster","tasks": ["task2"],"hosts": [{ "ip": "172.17.42.1", "sshPort": 32768, "username": "root", "description": "cassandra01", "password": "root"},{ "ip": "172.17.42.1", "sshPort": 32769, "description": "cassandra02", "username": "root", "interactiveAuth": true},{ "ip": "172.17.42.1", "sshPort": 32768, "description": "cassandra03", "username": "root", "password": "root"}] }]},"tasks": [{ "title": "task1", "maintainer": "Burkraith", "commands": [{"name": "ADD","sourcePath": "doc.go","destPath": "/tmp","description": "Raccoon.go to /tmp"},{"name": "RUN","description": "Removing htop","instruction": "sudo yum remove -y htop"},{"name": "ADD","sourcePath": "main.go","destPath": "/tmp","description": "copying raccoon.go to /tmp"} ]},{ "title": "task2", "maintainer": "Mario", "commands": [{"name": "RUN","description": "Removing htop","instruction": "sudo apt-get remove -y htop"} ]}] }`)
	r = bytes.NewReader(content)
	_, err = parseRequest(r)
	if err != nil {
		t.Fatal("An error must be thrown when passing a incorrect syntax like 'infras2tructure'")
	}
}
