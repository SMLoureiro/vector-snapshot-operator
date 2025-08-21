.PHONY: build run tidy
build:
\tgo build ./...

run:
\tgo run ./cmd/manager

tidy:
\tgo mod tidy
