# blogaggregator

A command-line RSS feed aggregator written in Go.

`blogaggregator` lets you register users, add RSS feeds, follow your favorites, and browse the latest posts from all the feeds you follow. A background aggregation loop periodically fetches every feed and saves new posts to a Postgres database.

---

## Prerequisites

You'll need two things installed on your machine:

- **Go** (1.21 or newer) — [download](https://go.dev/dl/)
- **PostgreSQL** (13 or newer) — [download](https://www.postgresql.org/download/)

Verify they're installed:

    go version      # go version go1.21.x ...
    psql --version  # psql (PostgreSQL) 15.x ...

### Setting up Postgres

Make sure the Postgres service is running and you have a database ready:

    # Create a database for the aggregator
    sudo -u postgres psql -c "CREATE DATABASE blog_aggregator;"

    # Set a password for the postgres user (if you haven't already)
    sudo -u postgres psql -c "ALTER USER postgres WITH PASSWORD 'yourpassword';"

Note the connection string — you'll need it in the configuration step:

    postgres://postgres:yourpassword@localhost:5432/blog_aggregator?sslmode=disable

---

## Installation

Install the `blogaggregator` CLI with `go install`:

    go install github.com/ihaliti03/blogaggregator@latest

This compiles the program and places the `blogaggregator` binary in `$(go env GOPATH)/bin` (usually `~/go/bin`).

Make sure that directory is on your `PATH`:

    export PATH="$PATH:$(go env GOPATH)/bin"

Verify the install:

    blogaggregator users

---

## Configuration

`blogaggregator` reads its configuration from `~/.gatorconfig.json`. Create it with your database connection string:

    cat > ~/.gatorconfig.json <<'EOF'
    {
      "db_url": "postgres://postgres:yourpassword@localhost:5432/blog_aggregator?sslmode=disable",
      "current_user_name": ""
    }
    EOF

Replace `yourpassword` with the password you set in the Postgres setup step.

### Apply the database schema

The project ships with SQL migrations in `sql/schema/`. Use [Goose](https://github.com/pressly/goose) to apply them:

    # Install goose (one-time)
    go install github.com/pressly/goose/v3/cmd/goose@latest

    # Apply all migrations
    goose -dir sql/schema postgres "postgres://postgres:yourpassword@localhost:5432/blog_aggregator?sslmode=disable" up

This creates the `users`, `feeds`, `feed_follows`, and `posts` tables.

---

## Usage

### Register a user

    blogaggregator register ibrahim

Creates a user and sets them as the current user.

### Add a feed

    blogaggregator addfeed "TechCrunch" "https://techcrunch.com/feed/"

Adds an RSS feed and automatically follows it for the current user.

### Follow another feed

    blogaggregator follow "https://news.ycombinator.com/rss"

Follows an existing feed (one that has already been added).

### List the feeds you follow

    blogaggregator following

### Start the aggregator

    blogaggregator agg 1m

Runs a long-lived background loop that fetches every feed once per interval (`1m` = one minute). New posts are saved to the database. Stop it with `Ctrl+C`.

**Tip:** run this in one terminal and use the other commands in a second terminal — the aggregator keeps the database up to date while you browse.

### Browse posts

    blogaggregator browse       # show the 2 most recent posts (default)
    blogaggregator browse 10    # show the 10 most recent posts

Displays the latest posts from all feeds you follow, newest first.

### Other commands

    blogaggregator login <username>   # switch to an existing user
    blogaggregator users              # list all registered users
    blogaggregator unfollow <url>     # stop following a feed

---

## Example workflow

    # 1. register and add feeds
    blogaggregator register ibrahim
    blogaggregator addfeed "TechCrunch" "https://techcrunch.com/feed/"
    blogaggregator addfeed "Hacker News" "https://news.ycombinator.com/rss"

    # 2. start the aggregator in the background (one terminal)
    blogaggregator agg 1m

    # 3. browse posts in another terminal
    blogaggregator browse 5

    # 4. follow additional feeds as you discover them
    blogaggregator follow "https://blog.boot.dev/index.xml"

---

## Development

To run from source instead of `go install`:

    git clone https://github.com/ihaliti03/blogaggregator.git
    cd blogaggregator
    go build -o blogaggregator .
    ./blogaggregator users

SQL queries are generated with [sqlc](https://sqlc.dev/):

    sqlc generate

Schema migrations use [Goose](https://github.com/pressly/goose):

    goose -dir sql/schema postgres "$DATABASE_URL" up
    goose -dir sql/schema postgres "$DATABASE_URL" down

---

## License

MIT