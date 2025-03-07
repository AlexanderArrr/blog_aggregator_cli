package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("handlerLogin() expected 1 argument, received %v", len(cmd.args))
	}

	userName := cmd.args[0]

	err := s.cfg.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("The Username %v has been set!\n", userName)
	return nil
}
