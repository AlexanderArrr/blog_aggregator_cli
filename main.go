package main

import (
	"context"
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
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("error while getting user: %v", err)
		}

		return handler(s, cmd, user)
	}
}
