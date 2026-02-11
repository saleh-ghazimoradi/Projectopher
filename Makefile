docker-up:
	docker compose up -d

docker-down:
	docker compose down

vet:
	go vet ./...

fmt:
	go fmt ./...

build:
	mkdir -p bin
	go build -o bin/projectopher

run: fmt vet
	go run . run
