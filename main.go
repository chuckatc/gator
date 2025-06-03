package main

import (
	"fmt"
	"log"

	"github.com/chuckatc/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	if err := cfg.SetUser("cat"); err != nil {
		log.Fatalln(err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(cfg.CurrentUserName)
	fmt.Println(cfg.DbURL)
}
