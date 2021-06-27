build:
	go build

test:
	go test ./...

install-hooks:
	git config core.hooksPath hooks
