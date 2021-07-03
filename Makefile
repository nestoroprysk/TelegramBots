.PHONY: deps 
deps:
	go mod tidy

.PHONY: build
build:
	go build

.PHONY: test 
test:
	go test ./...

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


.PHONY: sql-start 
sql-start:
	make sql-stop || true
	docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:5
	@echo "Starting sql..."
	@sleep 3
	@echo "Run the following command to connect:"
	@echo "  $$ mysql -P 3306 -u root -h 127.0.0.1 --password=root"

.PHONY: sql-stop 
sql-stop:
	docker kill mysql || true
	docker rm mysql
