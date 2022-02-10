ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build-api: 
	@echo "Building binary..." 
	@go build -o eKYC.o ./cmd/ekyc_api
	clear

run-api:
	@echo "Starting up docker..."
	@docker-compose up -d --remove-orphans
	make build-api
	@./eKYC.o --host=$(host) --db=$(db) --user=$(user) --pwd=$(pwd) --port=$(port) --server=$(server) --minio_server=$(minio_server) --minio_pwd=$(minio_pwd) --minio_user=$(minio_user) --redis_server=$(redis_server) --redis_pwd=$(redis_pwd)

build-daily-report: 
	@echo "Building binary..." 
	@go build -o daily_report.o ./cmd/daily_report
	clear

run-daily-report:
	@echo "Starting up docker..."
	@docker-compose up -d --remove-orphans
	make build-daily-report
	@./daily_report.o

build-monthly-report: 
	@echo "Building binary..." 
	@go build -o monthly_report.o ./cmd/monthly_report
	clear

run-monthly-report:
	@echo "Starting up docker..."
	@docker-compose up -d --remove-orphans
	make build-monthly-report
	@./monthly_report.o --client=$(client)

test:
	@go clean -testcache
	go test -v ./app/controllers/... ./app/middlewares
	
clean:
	@echo "Shutting down docker..."
	@docker-compose down
	@echo "Cleaning binaries..."
	@rm eKYC.o
	@go clean
	@go clean -testcache
	clear

