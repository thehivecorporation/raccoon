package raccoon

type Authentication struct {
	//Username to access remote host
	Username string `json:"username,omitempty"`

	//Password to access remote host
	Password string `json:"password,omitempty"`

	//Identity file for auth
	IdentityFile string `json:"identityFile,omitempty"`

	//Choose to enter user and password during Raccoon execution
	InteractiveAuth bool `json:"interactiveAuth,omitempty"`

	//SSHPort is a optional non standard port for SSH. We consider 22 as the
	//standard port so set this field in case you are using a different one in
	//any host
	SSHPort int `json:"sshPort,omitempty"`
}
