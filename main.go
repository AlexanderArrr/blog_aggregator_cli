package main

import (
	"fmt"
	"os/user"

	"github.com/AlexanderArrr/blog_aggregator_cli/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	currentUser, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}

	err = cfg.SetUser(currentUser.Username)
	if err != nil {
		fmt.Println(err)
	}

	cfg, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cfg)
}
