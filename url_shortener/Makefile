build:
	@go build -o bin/url_shortener

run: build
	@./bin/url_shortener $(ARGS)

test:
	@go test -v ./...
