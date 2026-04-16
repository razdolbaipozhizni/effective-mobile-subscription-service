.PHONY: help build run migrate-up migrate-down test swag docker-up docker-down

help: ## Показать все команды
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Собрать приложение
	go build -o bin/app ./cmd/app

run: ## Запустить локально (нужен запущенный Postgres)
	go run ./cmd/app

swag: ## Сгенерировать Swagger документацию
	swag init -g cmd/app/main.go --parseInternal --parseDependency

migrate-up: ## Применить миграции (нужен запущенный Postgres)
	go run github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path migrations -database "postgres://postgres:postgres@localhost:5432/subscriptions?sslmode=disable" up

migrate-down: ## Откатить миграции
	go run github.com/golang-migrate/migrate/v4/cmd/migrate@latest -path migrations -database "postgres://postgres:postgres@localhost:5432/subscriptions?sslmode=disable" down

docker-up: ## Запустить всё через docker compose
	docker compose up --build -d

docker-down: ## Остановить docker compose
	docker compose down -v

docker-logs: ## Посмотреть логи приложения
	docker compose logs -f app

test: ## Запустить тесты
	go test ./... -v