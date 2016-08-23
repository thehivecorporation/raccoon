package raccoon

type Command struct {
	Name        string
	Description string
}

func (c *Command) GetName() string {
	return c.Name
}
