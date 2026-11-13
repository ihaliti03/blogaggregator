package main

import "fmt"

func handlerLogin(s *state, cmd command) error{
	// don't check with cmd.arguments == nil but use len()
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("no argumments inputted")
	}

	if err := s.cfg.SetUser(cmd.arguments[0]);err != nil {
		return err
	}

	fmt.Printf("The user has been set to: %s\n", cmd.arguments[0])
	return nil
}