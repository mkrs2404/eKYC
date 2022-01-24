build: 
	@echo "Building binary..." 
	@go build -o eKYC.o main.go
	clear

run:
	@echo "Starting up docker..."
	@docker-compose up -d --remove-orphans
	make build
	./eKYC.o

test:
	go test -v ./api/controllers/...
	
clean:
	@echo "Shutting down docker..."
	@docker-compose down
	@echo "Cleaning binaries..."
	@rm eKYC.o
	@go clean
	clear

