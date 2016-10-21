package instructions

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/ssh"

	"fmt"
)

//IOWriterTester is a small utility to help writing unit tests where the output
//goes to a io.Writer.
type IOWriterTester struct {
	//Instance of *testing.T that you are using in your tests. IOWriterTester
	//will report here if it fails to find the expected string on the output
	T *testing.T

	//ExpectedString that you'll want to find in the io.Writer. It will parse
	//the output and use the T instance to report a not found string
	ExpectedString string

	//Set Log to true if you want to output to stdout the contents that passed
	//through the io.Writer
	Log bool
}

//Write implementation of io.Write interface that will capture contents written
//to this io.Writer and check their contents against the ExpectedString stored
//in the struct
func (t *IOWriterTester) Write(p []byte) (n int, err error) {
	content := string(p)
	if t.Log {
		t.T.Log(content)
	}
	if !strings.Contains(content, t.ExpectedString) {
		t.T.Errorf("String %s not found on io.Writer (Found '%s')\n", t.ExpectedString, content)
	}
	return len(p), nil
}

func passOrError(err error, f func(args ...interface{})) {
	if err != nil {
		f(err)
	}
}

//Session must be closed manually
func createSessionPrototype(hostIp string, port int) (*ssh.Session, error) {
	clientConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("root"),
		},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", hostIp, port), clientConfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to dial: " + err.Error())
	}
	session, err := client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: " + err.Error())
	}
	return session, nil
}

func createTestFileWithContent(c string) (*os.File, error) {
	os.MkdirAll("/tmp/raccoon", 0777)

	f, err := ioutil.TempFile("/tmp/raccoon/", "raccoon")
	if err != nil {
		return nil, err
	}

	_, err = f.WriteString(c)
	if err != nil {
		return nil, err
	}

	if err = f.Sync(); err != nil {
		return nil, err
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}

	return f, nil
}
