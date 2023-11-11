build:
	@go build -o bin/cli_task_manager

run: build
	@./bin/cli_task_manager $(ARGS)

install:
	@go install .

test:
	@go test ./...
