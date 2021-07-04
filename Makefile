.PHONY: deps 
deps:
	go mod tidy

.PHONY: build
build:
	go build

.PHONY: test
test:
	go test ./...

.PHONY: test-ci
test-ci: up
	docker-compose run test

.PHONY: cover
cover:
	go test ./... -coverprofile cover.out
	go tool cover -html=cover.out -o cover.html
	open cover.html

.PHONY: doc
doc:
	@echo "Run the following command to open the documentation:"
	@echo "  $$ python -m webbrowser http://localhost:6060/pkg/github.com/nestoroprysk/TelegramBots/?m=all"
	godoc -http=:6060

.PHONY: install-hooks 
install-hooks:
	git config core.hooksPath hooks

host := $(if ${CI},"mysql","127.0.0.1")

.PHONY: up 
up:
	docker network create cloudbuild || true
	docker-compose up -d --remove-orphans mysql || true
	./retry mysql -P 3306 -u root -h ${host} --password=root -e "select 1;"
	@echo "Success!"

.PHONY: down 
down:
	docker-compose down
