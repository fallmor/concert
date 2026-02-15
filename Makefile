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

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-deploy: nix-docker-load docker-up
	@echo "Deployed with Nix pkgs images!"

# Use system git for nix-portable
export NP_GIT := $(shell which git 2>/dev/null || echo "")
NIX_CMD := $(shell which nix-portable >/dev/null 2>&1 && echo "nix-portable nix" || echo "nix")

# Nix Docker build 
nix-docker-build:
	@echo "Cleaning up /homeless-shelter if it exists (nix-portable issue)"
	@echo "Building server image with Nix pkgs"
	$(NIX_CMD) build '.#docker' --impure --no-sandbox -o server-image
	@echo "Building worker image with Nix pkgs"
	@sudo rm -rf /homeless-shelter 2>/dev/null || true
	$(NIX_CMD) build '.#dockerWorker' --impure --no-sandbox -o worker-image

nix-docker-load: nix-docker-build
	@echo "Loading server image with Nix pkgs"
	docker load < server-image
	docker load < worker-image
	@rm -f server-image worker-image
	@echo "Docker images loaded successfully!"

nix-docker: nix-docker-load
	@echo "Server and worker images built and loaded!"
	@echo "Run 'make docker-up' to start all services"

docker-up-nix: nix-docker-load docker-up
	@echo "All services started with Nix-built images!"

# Nix development shell
nix-shell:
	$(NIX_CMD) develop

# Build individual components with Nix
nix-server:
	$(NIX_CMD) build .#server

nix-frontend:
	$(NIX_CMD) build .#frontend

# Pre-build frontend (useful to avoid npm install in Nix)
prebuild-frontend:
	@echo "Pre-building frontend to avoid npm install in Nix..."
	cd frontend && npm install --legacy-peer-deps && npm run build
	@echo "Frontend pre-built! Now run: make nix-docker-build"