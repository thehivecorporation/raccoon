package connection

import "testing"

//func TestGetSession(t *testing.T) {
//	n := Node{
//		IP:"192.168.33.10",
//		Username:"vagrant",
//		Password:"vagrant",
//	}
//
//	session, err := n.GetSession()
//
//	if err != nil || session == nil{
//		t.Fatal(err)
//	}
//}

func TestGetOutput(t *testing.T){
	n := Node{
		IP:"192.168.33.10",
		Username:"vagrant",
		Password:"vagrant",
	}

	n.GetOutput()
}