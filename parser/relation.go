package parser

import "github.com/thehivecorporation/raccoon"

type Relation struct {
	Generic
}

func (rp *Relation) Prepare(i *raccoon.Infrastructure, rl *raccoon.RelationList) error {
	for _, relation := range *rl {
		for k, cluster := range i.Infrastructure {
			if cluster.Name == relation.ClusterName {
				i.Infrastructure[k].TasksToExecute = relation.TaskList
			}
		}
	}

	return nil
}
