package parser

import "fmt"

const (
	PARSING_ERROR  string = "Error found when parsing infrastructure file"
	BLANK_IP       string = "Host IP can't be blank"
	BLANK_PASSWORD string = "Host password is blank on host '%s'. If no password is specified you " +
		"must use an identity file or an interactive authentication method"
	BLANK_USERNAME string = "Host username is blank on host '%s'"
	NO_TASKS       string = "You haven't specified any task. Specify at least one as an string " +
		"array on cluster '%s'"
	NO_HOSTS        string = "No hosts were found on cluster '%s' for commands '%s'"
	NO_CLUSTER_NAME string = "infrastructure name can't be blank"
	NO_CLUSTER      string = "No cluster was found on infrastructure file"
	JSON_ERROR      string = "Error parsing JSON: %s\n"
)

type ParseError struct {
	errorCode string
	msg       string
	extra     []string
}

func (i *ParseError) Error() string {
	return fmt.Sprintf(i.msg, i.extra)
}

func parseErrorFactory(errorCode string, extra ...string) error {
	err := ParseError{
		errorCode: errorCode,
		msg:       errorCode,
		extra:     extra,
	}

	return &err
}
