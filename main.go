package main

import (
	"database/sql"
	"log"
	"os"
	"path"

	"github.com/chuckatc/gator/internal/config"
	"github.com/chuckatc/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalln(err)
	}
	dbQueries := database.New(db)

	s := state{
		db:  dbQueries,
		cfg: &cfg,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <command> [args ...]", path.Base(os.Args[0]))
	}
	cmdName := os.Args[1]
	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	cmd := command{name: cmdName, args: args}
	if err := cmds.run(&s, cmd); err != nil {
		log.Fatalln(err)
	}
}
