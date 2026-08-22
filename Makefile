APP_NAME=go-repair-center

.PHONY: test size frontend-build docker-up

test:
	go test ./...

size:
	powershell -ExecutionPolicy Bypass -File scripts/verify-go-size.ps1

frontend-build:
	cd frontend && npm run build

docker-up:
	docker compose -f docker-compose.yml up --build
