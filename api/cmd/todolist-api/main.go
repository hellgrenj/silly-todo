package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hellgrenj/silly-todo/pkg/todo"

	"github.com/hellgrenj/silly-todo/pkg/storage"

	"github.com/hellgrenj/silly-todo/pkg/http/rest"
)

func main() {
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	// first init db
	db := storage.NewDatabase(1)
	// then construct todo service (passing in db dependency)
	todo := todo.NewService(db) // db qualifies as a todo repository .. meets interface contract
	// then constructing rest server (passing in service dependency)
	s := rest.NewServer(todo)
	// then start rest service
	log.Fatal(s.StartRestServer(":8080"))
}
