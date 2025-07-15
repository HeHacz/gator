package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/hehacz/gator/internal/config"
	"github.com/hehacz/gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db   *database.Queries
	conf *config.Config
}

func main() {
	conf, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
	}
	db, err := sql.Open("postgres", conf.DB_url)
	if err != nil {
		log.Fatalf("error couldnt connect to the database: %v", err)
	}
	instanceState := &state{
		db:   database.New(db),
		conf: &conf,
	}
	cmds := commands{
		avaiableCommands: make(map[string]func(*state, command) error),
	}

	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)

	if len(os.Args) < 2 {
		log.Fatal("At least one argument is required")
		return
	}
	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	if err := cmds.run(instanceState, command{name: cmdName, args: cmdArgs}); err != nil {
		log.Fatal(err)
	}
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config: %+v\n", cfg)

}
