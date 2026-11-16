package main

import (
	//	"fmt"
	"database/sql"
	"log"
	"os"

	"github.com/ihaliti03/blogaggregator/internal/config"
	"github.com/ihaliti03/blogaggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	dbQuery := database.New(db)
	s := &state{
		cfg: &cfg,
		db:  dbQuery, //Now every handler can call s.db.GetUser(ctx, name) and s.db.CreateUser(ctx, params)
	}

	cmds := &commands{handlers: make(map[string]func(*state, command) error)}

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
