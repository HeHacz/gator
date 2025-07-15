package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	avaiableCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	fun, ok := c.avaiableCommands[cmd.name]
	if !ok {
		return fmt.Errorf("command not found")
	}
	return fun(s, cmd)
}

func (c *commands) register(name string, fun func(*state, command) error) {
	c.avaiableCommands[name] = fun
}
