package main

import (
//	"fmt"
	"log"
	"os"

	"github.com/ihaliti03/blogaggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	s := &state{cfg: &cfg}
    cmds := &commands{handlers: make(map[string]func(*state, command) error)}

    cmds.register("login", handlerLogin)

	if len(os.Args) < 2 {
        log.Fatal("no command provided")
    }
    cmd := command{
        name:      os.Args[1],
        arguments: os.Args[2:],
    }
	if err := cmds.run(s, cmd); err != nil {
        log.Fatal(err)
    }

	
	//fmt.Printf("%+v\n", updatedConfig)
}
