.PHONY: build run test sqlc goose-up goose-down clean

# Default database URL — override with: make goose-up DATABASE_URL=...
DATABASE_URL ?= postgres://postgres:yourpassword@localhost:5432/blog_aggregator?sslmode=disable

build:
	go build -o blogaggregator .

run: build
	./blogaggregator

test:
	go test ./...

sqlc:
	sqlc generate

goose-up:
	goose -dir sql/schema postgres "$(DATABASE_URL)" up

goose-down:
	goose -dir sql/schema postgres "$(DATABASE_URL)" down

clean:
	rm -f blogaggregator
	go clean -cache