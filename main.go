package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/chuckatc/gator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	s := state{
		cfg: &cfg,
	}

	cmds := commands{
		handlers: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <command> [args ...]", path.Base(os.Args[0]))
	}
	cmdName := os.Args[1]
	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	cmds.run(&s, command{name: cmdName, args: args})

	fmt.Println("", cmdName, args)
	fmt.Println("", s, cmds)

	fmt.Println(cfg.CurrentUserName)
	fmt.Println(cfg.DbURL)
}
