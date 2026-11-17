package handler

import "fmt"

type Commands struct {
	Handlers map[string]HandlerFunc
}

func NewCommands() *Commands {
	return &Commands{Handlers: make(map[string]HandlerFunc)}
}

func (c *Commands) Register(name string, f HandlerFunc) {
	c.Handlers[name] = f
}

func (c *Commands) Run(s *State, cmd Command) error {
	h, ok := c.Handlers[cmd.Name]
	if !ok {
		return fmt.Errorf("unknown command: %s", cmd.Name)
	}
	return h(s, cmd)
}
