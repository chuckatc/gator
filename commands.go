package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/chuckatc/gator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

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
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected username, got %v", cmd.args)
	}
	userName := cmd.args[0]

	_, err := s.db.GetUser(context.Background(), userName)
	if errors.Is(err, sql.ErrNoRows) {
		log.Fatalf("User not found: %s", userName)
	}
	if err != nil {
		return err
	}

	if err := s.cfg.SetUser(userName); err != nil {
		return err
	}
	fmt.Printf("User set to: %s\n", userName)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("expected username, got: %v", cmd.args)
	}
	userName := cmd.args[0]

	createUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      userName,
	}
	dbUser, err := s.db.CreateUser(context.Background(), createUserParams)
	if isUniqueViolation(err) {
		log.Fatalf("User already exists: %s", userName)
	}
	if err != nil {
		return err
	}

	if err := s.cfg.SetUser(userName); err != nil {
		return err
	}
	fmt.Printf("User created: %s\n", dbUser.Name)
	fmt.Printf("ID: %s\n", dbUser.ID)
	fmt.Printf("Created At: %s\n", dbUser.CreatedAt)
	fmt.Printf("Updated At: %s\n", dbUser.UpdatedAt)

	return nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == "23505"
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) != 0 {
		return fmt.Errorf("unexpected arguments to reset: %v", cmd.args)
	}

	if err := s.db.TruncateUsers(context.Background()); err != nil {
		return err
	}

	fmt.Println("Users reset successfully")
	return nil
}
