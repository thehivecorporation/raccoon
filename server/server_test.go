package server

import (
	"bytes"
	"testing"
)

func TestParseRequest(t *testing.T) {
	content := []byte("{\"zombiebook\":[{\"chapter_title\": \"chapter1\", \"maintainer\": \"Burkraith\", \"instructions\": [] }, {\"chapter_title\": \"chapter2\", \"maintainer\": \"Mario\", \"instructions\": [] } ], \"mansion\":{\"name\":\"A name\", \"rooms\":[{\"name\":\"some room\", \"chapter\":\"chapter1\", \"hosts\":[{\"ip\":\"192.168.1.44\", \"username\":\"vagrant\", \"password\":\"vagrant\"} ] }, {\"name\":\"some room2\", \"chapter\":\"chapter2\", \"hosts\":[{\"ip\":\"192.168.1.45\", \"username\":\"vagrant\", \"password\":\"vagrant\"} ] } ] } }")
	r := bytes.NewReader(content)
	err, _ := parseRequest(r)
	if err != nil {
		t.Fatal(err)
	}

	content = []byte("{wrong:\"syntax\",}")
	r = bytes.NewReader(content)
	err, _ = parseRequest(content)
	if err == nil {
		t.Fatal("An error must be thrown when using a wrong JSON syntax")
	}

	content = []byte("{\"zombiebook\":[{\"chapter_title\": \"chapter1\", \"maintainer\": \"Burkraith\", \"instructions\": [] }, {\"chapter_title\": \"chapter2\", \"maintainer\": \"Mario\", \"instructions\": [{\"name\": \"RUN\", \"description\": \"Install wget\", \"inst2ruction\": \"sudo yum install -y wget\"} ] } ], \"mansion\":{\"name\":\"A name\", \"rooms\":[{\"name\":\"some room\", \"chapter\":\"chapter1\", \"hosts\":[{\"ip\":\"192.168.1.44\", \"username\":\"vagrant\", \"password\":\"vagrant\"} ] }, {\"name\":\"some room2\", \"chapter\":\"chapter2\", \"hosts\":[{\"ip\":\"192.168.1.45\", \"username\":\"vagrant\", \"password\":\"vagrant\"} ] } ] } }")
	r = bytes.NewReader(content)
	err, _ = parseRequest(content)
	if err == nil {
		t.Fatal("An error must be thrown when passing a incorrect syntax like 'inst2ruction'")
	}
}
