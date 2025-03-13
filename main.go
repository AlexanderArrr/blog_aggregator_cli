package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/AlexanderArrr/blog_aggregator_cli/internal/config"
	"github.com/AlexanderArrr/blog_aggregator_cli/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	// Config and State initialization
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	programState := &state{
		cfg: &cfg,
	}

	// Registering commands
	cmds := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)

	// Database connection
	db, err := sql.Open("postgres", cfg.DbUrl)
	if err != nil {
		fmt.Println("Error while loading database")
	}

	programState.db = database.New(db)

	if len(os.Args) < 2 {
		log.Fatal("Usage: gator <command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	cmd := command{
		name: cmdName,
		args: cmdArgs,
	}

	err = cmds.run(programState, cmd)
	if err != nil {
		log.Fatalf("error while executing command: %v", err)
	}

}
