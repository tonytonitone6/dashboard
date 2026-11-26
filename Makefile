.PHONY: setup dev api web build test deploy

HOST ?= stan@34.55.165.112
SITE ?= https://stanmao.com

setup:
	cd api && go mod tidy
	cd web && pnpm install

# Make runs this rule only when the file does not exist.
api/.env:
	@echo "api/.env is missing. Copy api/.env.example to api/.env and fill in every value." >&2; exit 1

api: api/.env
	cd api && set -a && . ./.env && set +a && go run ./cmd/server

web:
	cd web && pnpm dev

dev: api/.env
	$(MAKE) -j2 api web

build:
	cd api && go build -o bin/server ./cmd/server
	cd web && pnpm build

test:
	cd api && go vet ./... && go test ./...
	cd web && pnpm lint

# The new binary checks the VM's .env before anything in use is replaced, so a
# bad .env stops the deploy and leaves the running API alone.
deploy: test
	ssh $(HOST) 'test -f apps/dashboard/.env || { echo "apps/dashboard/.env is missing on this VM. Create it with the variables in api/.env.example, then deploy again." >&2; exit 1; }'
	cd web && pnpm build
	cd api && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o bin/linux-amd64/server ./cmd/server
	ssh $(HOST) 'mkdir -p apps/dashboard/bin apps/nginx/html/dashboard'
	rsync -az deploy/compose.yaml $(HOST):apps/dashboard/compose.yaml
	rsync -az api/bin/linux-amd64/server $(HOST):apps/dashboard/bin/server.new
	ssh $(HOST) 'cd apps/dashboard && docker compose run --rm --no-deps -T api /app/server.new check-config'
	rsync -az --delete-after --exclude .DS_Store web/dist/ $(HOST):apps/nginx/html/dashboard/
	ssh $(HOST) 'cd apps/dashboard && mv bin/server.new bin/server && docker compose up -d --force-recreate'
	curl -fsS --retry 5 --retry-all-errors --retry-delay 2 $(SITE)/api/health || { ssh $(HOST) 'docker logs --tail 20 dashboard-api'; exit 1; }
