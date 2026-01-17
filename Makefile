.PHONY: build test run clean frontend-build server-build docker-build docker-up docker-down docker-logs

build: frontend-build server-build

frontend-build:
	cd frontend && npm run build

server-build:
	go build -o server cmd/server/main.go

test:
	go test ./...
fmt:
	go fmt ./...

run:
	go run cmd/server/main.go

clean:
	rm -f server
	rm -rf frontend/dist
	rm -rf static/assets
	rm -f static/index.html

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down


docker-deploy: docker-build docker-up