.PHONY: help build start stop clean logs test dev prod

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build all Docker images
	docker-compose build

build-backend: ## Build backend Docker image
	docker-compose build backend

build-frontend: ## Build frontend Docker image
	docker-compose build frontend

start: ## Start all services
	docker-compose up -d

stop: ## Stop all services
	docker-compose down

restart: ## Restart all services
	docker-compose restart

clean: ## Stop and remove all containers, networks, and volumes
	docker-compose down -v

logs: ## View logs from all services
	docker-compose logs -f

logs-backend: ## View backend logs
	docker-compose logs -f backend

logs-frontend: ## View frontend logs
	docker-compose logs -f frontend

logs-db: ## View database logs
	docker-compose logs -f postgres

ps: ## Show running containers
	docker-compose ps

dev: ## Start development environment
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up

dev-build: ## Build and start development environment
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build

prod: ## Start production environment
	docker-compose up -d

shell-backend: ## Open shell in backend container
	docker-compose exec backend sh

shell-frontend: ## Open shell in frontend container
	docker-compose exec frontend sh

shell-db: ## Open PostgreSQL shell
	docker-compose exec postgres psql -U gnucash -d gnucash

backup-db: ## Backup database
	@mkdir -p backups
	docker-compose exec -T postgres pg_dump -U gnucash gnucash > backups/backup-$$(date +%Y%m%d-%H%M%S).sql
	@echo "Database backed up to backups/backup-$$(date +%Y%m%d-%H%M%S).sql"

restore-db: ## Restore database (usage: make restore-db FILE=backup.sql)
	@if [ -z "$(FILE)" ]; then \
		echo "Error: FILE parameter required. Usage: make restore-db FILE=backup.sql"; \
		exit 1; \
	fi
	docker-compose exec -T postgres psql -U gnucash -d gnucash < $(FILE)
	@echo "Database restored from $(FILE)"

test-backend: ## Run backend tests
	cd backend && go test ./...

test-frontend: ## Run frontend tests
	cd frontend && npm test

health: ## Check health of all services
	@echo "Backend health:"
	@curl -s http://localhost:8080/health | jq || echo "Backend not responding"
	@echo "\nFrontend health:"
	@curl -s -o /dev/null -w "%{http_code}" http://localhost:3000 && echo " OK" || echo " Failed"

prune: ## Remove all unused Docker resources
	docker system prune -af --volumes

install-backend: ## Install backend dependencies
	cd backend && go mod download

install-frontend: ## Install frontend dependencies
	cd frontend && npm install

run-backend-local: ## Run backend locally (without Docker)
	cd backend && go run ./cmd/server

run-frontend-local: ## Run frontend locally (without Docker)
	cd frontend && npm start
