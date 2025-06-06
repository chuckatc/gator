package main

import "fmt"

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	callback, found := c.handlers[cmd.name]
	if !found {
		return fmt.Errorf("command not found: %s", cmd.name)
	}

	err := callback(s, cmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	fmt.Println(len(cmd.args))
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected username, got %v", cmd.args)
	}

	userName := cmd.args[0]
	if err := s.cfg.SetUser(userName); err != nil {
		return err
	}
	fmt.Printf("User set to: %s\n", userName)

	return nil
}
