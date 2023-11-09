buildweb:
	@go build -C cmd/web -o ../../bin/cyoaweb

runweb: buildweb
	@./bin/cyoaweb $(ARGS)

buildcli:
	@go build -C cmd/cli -o ../../bin/cyoacli

runcli: buildcli
	@./bin/cyoacli $(ARGS)

build: buildweb

run: runweb

test:
	@go test ./...
