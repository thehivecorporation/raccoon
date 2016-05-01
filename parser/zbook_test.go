package parser

import (
	"io/ioutil"
	"testing"
)

func TestReadZbookFile(t *testing.T) {
	//Non existing filepath
	filePath := "/tmp/i-do-not-exist"

	_, err := readZbookFile(filePath)
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

	_, err = readZbookFile(filePath)
	if err == nil {
		t.Fatal("After passing a json with wrong structure an error must be thrown")
	}

	//Show error if you use "go test -v"
	t.Log(err)

	//Pass a correct JSON
	filePath = "../examples/exampleBook.json"

	_, err = readZbookFile(filePath)
	if err != nil {
		t.Fatalf("JSON file was correct. No error must be here: %s\n", err.Error())
	}

	//Pass a syntactically correct zbook with no content
	f, _ = ioutil.TempFile("/tmp", "no-content")
	f.WriteString("[{\"chapter_title\":\"chapter1\",\"maintainer\": \"Burkraith\",\"instructions\": []}]")
	f.Close()
	filePath = f.Name()

	_, err = readZbookFile(filePath)
	if err != nil {
		t.Fatal("No instructions were provided but syntax was correct")
	}

	//Pass a syntactically correct zbook with incorrect keys or no present
	f, _ = ioutil.TempFile("/tmp", "incorrect-keys")
	f.WriteString("[{\"chapter_title\": \"chapter1\", \"maintainer\": \"Burkraith\", \"instructions\": [{\"na2me\": \"RUN\", \"description\": \"Install htop\", \"instruction\": \"sudo yum install -y htop\"} ] } ]")
	f.Close()
	filePath = f.Name()

	_, err = readZbookFile(filePath)
	if err == nil {
		t.Fatal("No instructions were provided but syntax was correct")
	}

	//Pass a syntactically correct zbook with incorrect keys or no present
	f, _ = ioutil.TempFile("/tmp", "incorrect-keys")
	f.WriteString("[{\"chapter_title\": \"chapter1\", \"maintainer\": \"Burkraith\", \"instructions\": [{\"name\": \"RUN\", \"description\": \"Install htop\", \"instruct2ion\": \"sudo yum install -y htop\"} ] } ]")
	f.Close()
	filePath = f.Name()

	_, err = readZbookFile(filePath)
	if err == nil {
		t.Fatal("No instructions were provided but syntax was correct")
	}

	//Pass a syntactically correct zbook with incorrect keys or no present
	f, _ = ioutil.TempFile("/tmp", "incorrect-keys")
	f.WriteString("[{\"chapter_title\": \"chapter1\", \"maintainer\": \"Burkraith\", \"instructions\": [{\"name\": \"RUN\", \"description\": \"Install htop\", \"instruction\": \"sudo yum install -y htop\"}, {\"name\": \"ADD\", \"description\": \"Copying conf file\", \"source2Path\": \"/tmp/asdfad\", \"destPath\": \"/tmp/folder\"} ] } ]")
	f.Close()
	filePath = f.Name()

	_, err = readZbookFile(filePath)
	if err == nil {
		t.Fatal("No instructions were provided but syntax was correct")
	}

	//Pass a syntactically correct zbook with incorrect keys or no present
	f, _ = ioutil.TempFile("/tmp", "incorrect-keys")
	f.WriteString("[{\"chapter_title\": \"chapter1\", \"maintainer\": \"Burkraith\", \"instructions\": [{\"name\": \"RUN\", \"description\": \"Install htop\", \"instruction\": \"sudo yum install -y htop\"}, {\"name\": \"ADD\", \"description\": \"Copying conf file\", \"sourcePath\": \"/tmp/asdfad\", \"destP3ath\": \"/tmp/folder\"} ] } ]")
	f.Close()
	filePath = f.Name()

	_, err = readZbookFile(filePath)
	if err == nil {
		t.Fatal("No instructions were provided but syntax was correct")
	}
}
