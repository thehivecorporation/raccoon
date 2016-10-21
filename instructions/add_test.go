package instructions

import (
	"testing"
	"time"

	"github.com/thehivecorporation/raccoon"
)

var targetHostIP = "127.0.0.1"
var targetUsername = "root"
var targetPassword = "root"

func Test_Integration_ExecuteAddInstruction(t *testing.T) {
	host := getPrototypeHost()

	testContents := "Hello raccoon2\n"
	f, err := createTestFileWithContent(testContents)
	passOrError(err, t.Fatal)

	defer f.Close()

	add := ADD{
		Command: raccoon.Command{
			Name:        "ADD",
			Description: "Add description",
		},
		DestPath:   "/tmp/raccoon",
		SourcePath: f.Name(),
	}

	err = host.InitializeNode()
	passOrError(err, t.Fatal)

	add.Execute(host)
	host.CloseNode()

	time.Sleep(1 * time.Second)

	testSession, err := createSessionPrototype(targetHostIP, 22)
	passOrError(err, t.Fatal)
	defer testSession.Close()

	testWriter := &IOWriterTester{
		T:              t,
		ExpectedString: testContents,
		Log:            true,
	}
	testSession.Stdout = testWriter
	testSession.Stderr = testWriter

	if err := testSession.Run("cat /tmp/raccoon && rm /tmp/raccoon"); err != nil {
		t.Error(err)
	}
}

func getPrototypeHost() raccoon.Host {
	host := raccoon.Host{
		IP: targetHostIP,
		Authentication: raccoon.Authentication{
			Username: targetUsername,
			Password: targetPassword,
		},
	}

	return host
}

func Test_Integration_CopyFileToHost(t *testing.T) {
	host := getPrototypeHost()

	add := ADD{
		Command: raccoon.Command{
			Name:        "ADD",
			Description: "Adding a file to destination target",
		},
		DestPath: "/tmp/raccoon",
	}

	err := host.InitializeNode()
	passOrError(err, t.Fatal)

	//Prepare: Create a temp file to pass in with the contents "Hello Raccoon"
	//(without quotes)
	fileTestcontent := "Hello raccoon\n"
	f, err := createTestFileWithContent(fileTestcontent)
	passOrError(err, t.Fatal)

	add.SourcePath = f.Name()
	defer f.Close()

	//Prepare: Create a session with the container
	time.Sleep(1 * time.Second) //Give the container some time to launch
	session, err := host.GetSession()
	passOrError(err, t.Fatal)

	defer session.Close()

	//Prepare: Send the file with the information we have until now to /tmp
	err = copyToSession(session, add.DestPath, f)
	passOrError(err, t.Error)

	//Test: Check that the file exists in the target machine on path /tmp. The
	//file must have the same name and content in the target machine
	testSession, err := createSessionPrototype(host.IP, host.SSHPort)
	passOrError(err, t.Fatal)

	defer testSession.Close()

	testWriter := &IOWriterTester{
		T:              t,
		ExpectedString: fileTestcontent,
		Log:            true,
	}
	testSession.Stdout = testWriter
	testSession.Stderr = testWriter

	if err := testSession.Run("cat /tmp/raccoon && rm /tmp/raccoon"); err != nil {
		t.Error(err)
	}
}
