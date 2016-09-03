package raccoon

type Relation struct {
	ClusterName string   `json:"clusterName"`
	TaskList    []string `json:"tasks"`
}

type RelationList []Relation

