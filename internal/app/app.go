package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/ihaliti03/blogaggregator/internal/config"
	"github.com/ihaliti03/blogaggregator/internal/database"
	"github.com/ihaliti03/blogaggregator/internal/handler"
	"github.com/ihaliti03/blogaggregator/internal/middleware"
	_ "github.com/lib/pq"
)

func Run() {
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
	s := &handler.State{
		Cfg: &cfg,
		Db:  dbQuery, //Now every handler can call s.db.GetUser(ctx, name) and s.db.CreateUser(ctx, params)
	}

	cmds := &handler.Commands{Handlers: make(map[string]handler.HandlerFunc)}

	cmds.Register("login", handler.HandleLogin)
	cmds.Register("register", handler.HandleRegister)
	cmds.Register("reset", handler.HandleReset)
	cmds.Register("users", handler.HandleUsers)
	cmds.Register("agg", handler.HandleAgg)
	cmds.Register("addfeed", middleware.LoggedIn(handler.HandleAddFeed))
	cmds.Register("feeds", handler.HandleFeeds)
	cmds.Register("follow", middleware.LoggedIn(handler.HandleFollow))
	cmds.Register("following", middleware.LoggedIn(handler.HandleFollowing))
	cmds.Register("unfollow", middleware.LoggedIn(handler.HandleUnfollow))
	cmds.Register("browse", middleware.LoggedIn(handler.HandleBrowse))
	cmds.Register("deletefeed", middleware.LoggedIn(handler.HandleDeleteFeed))

	if len(os.Args) < 2 {
		log.Fatal("no command provided")
	}
	cmd := handler.Command{
		Name:      os.Args[1],
		Arguments: os.Args[2:],
	}
	if err := cmds.Run(s, cmd); err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%+v\n", updatedConfig)
}
