package raccoon

type Request struct {
	CommandsFile   []RawTask   `json:"commandsList"`
	Infrastructure Infrastructure `json:"infrastructure"`
}
