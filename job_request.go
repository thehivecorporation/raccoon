package raccoon

//JobRequest is a job defined by the user in a single JSON object
type JobRequest struct {
	TaskList       *[]Task         `json:"tasks"`
	Infrastructure *Infrastructure `json:"infrastructure"`
}
