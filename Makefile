MIGRATE = migrate
MIGRATIONS_DIR = db/migrations
DB_URL = "mysql://root:fadel123@tcp(localhost:9307)/isp_management?charset=utf8mb4&parseTime=True&loc=Local"

## Create new migration
create-migration:
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) $(name)

## Run migrations
migrate-up:
	$(MIGRATE) -database $(DB_URL) -path $(MIGRATIONS_DIR) up

## Rollback 1 step
migrate-down:
	$(MIGRATE) -database $(DB_URL) -path $(MIGRATIONS_DIR) down 1

## Drop all migrations
migrate-drop:
	$(MIGRATE) -database $(DB_URL) -path $(MIGRATIONS_DIR) drop -f

### Start docker development
mysql:
	docker container start mysql-container1

app:
	go run cmd/web/main.go

# --- DOCKER DEPLOYMENT ---

## Build and run all services in docker (detached mode)
docker-up:
	docker compose up --build -d

## Stop all running services in docker
docker-down:
	docker compose down

## Restart docker containers
docker-restart:
	docker compose restart

## Show live logs of all services
docker-logs:
	docker compose logs -f

## Show status of all services
docker-status:
	docker compose ps

# --- QA DEPLOYMENT ---
QA_USER = root
QA_HOST = 172.16.23.70
QA_PATH = /app/isp-management

## Deploy and launch on QA Server (172.16.23.70)
deploy-qa:
	@echo "Preparing QA directory..."
	ssh $(QA_USER)@$(QA_HOST) "mkdir -p $(QA_PATH)"
	@echo "Uploading files to QA Server..."
	scp -r Dockerfile docker-compose.yml .env config.json go.mod go.sum Makefile cmd db internal frontend $(QA_USER)@$(QA_HOST):$(QA_PATH)/
	@echo "Building and starting containers on QA Server..."
	ssh $(QA_USER)@$(QA_HOST) "cd $(QA_PATH) && docker compose down && docker compose up --build -d"
	@echo "Deployment to QA Server completed!"