package connection

//Node is a remote machine (virtual or physical) that we will execute our
//instructions on. Fields are self descriptive. AuthFilePath corresponds to the
//path of the private key that could give access to a remote machine (TODO)
type Node struct {
	IP           string
	Username     string
	Password     string
	AuthFilePath string
}
