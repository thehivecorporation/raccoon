package job

import "github.com/thehivecorporation/raccoon/connection"

type Job struct {
	connection.Cluster
	Recipe
}
