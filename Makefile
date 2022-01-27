ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: 
	@echo "Building binary..." 
	@go build -o eKYC.o
	clear

run:
	@echo "Starting up docker..."
	@docker-compose up -d --remove-orphans
	make build
	@./eKYC.o --host=$(host) --db=$(db) --user=$(user) --pwd=$(pwd) --port=$(port) --server=$(server) --minio_server=$(minio_server) --minio_pwd=$(minio_pwd) --minio_user=$(minio_user)

test:
	go test -v ./api/controllers/...
	
clean:
	@echo "Shutting down docker..."
	@docker-compose down
	@echo "Cleaning binaries..."
	@rm eKYC.o
	@go clean
	clear

