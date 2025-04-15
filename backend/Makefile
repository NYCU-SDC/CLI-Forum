GREEN = \033[0;32m
BLUE = \033[0;34m
RED = \033[0;31m
NC = \033[0m

all: gen build test
	@echo -e "==> $(BLUE)All tasks completed successfully$(NC)"

run: gen
	@echo -e ":: $(GREEN)Starting backend...$(NC)"
	@go build -o bin/backend cmd/backend/main.go && \
		DEBUG=true bin/backend \
		&& echo -e "==> $(BLUE)Successfully shutdown backend$(NC)" \
		|| (echo -e "  -> $(RED)Backend failed to start$(NC)" && exit 1)

observe: gen
	@echo -e ":: $(GREEN)Starting observability stack...$(NC)"
	@docker compose -f observe/docker-compose.yml up --quiet-pull -d \
		&& echo -e "  -> Observability stack started" \
		|| (echo -e "  -> $(RED)Failed to start observability stack$(NC)" && exit 1)

	$(call wait_healthy,observe-loki-1)

	@echo -e ":: $(GREEN)Streaming logs to Vector (local binary) with OTEL enabled...$(NC)"
	@echo -e "  -> If start successfully, you can open http://localhost:3000/explore?schemaVersion=1&panes=%7B%2251t%22%3A%7B%22datasource%22%3A%22P8E80F9AEF21F6940%22%2C%22queries%22%3A%5B%7B%22refId%22%3A%22A%22%2C%22expr%22%3A%22%7Bapp%3D%5C%22cli-forum-dev%5C%22%7D+%7C+json+%7C+line_format+%5C%22%5B%7B%7B+.level+%7D%7D%5D+%7B%7B+.msg+%7D%7D+%28%7B%7B+.caller+%7D%7D%29%5C%22%22%2C%22queryType%22%3A%22range%22%2C%22datasource%22%3A%7B%22type%22%3A%22loki%22%2C%22uid%22%3A%22P8E80F9AEF21F6940%22%7D%2C%22editorMode%22%3A%22code%22%2C%22direction%22%3A%22backward%22%7D%5D%2C%22range%22%3A%7B%22from%22%3A%22now-24h%22%2C%22to%22%3A%22now%22%7D%2C%22panelsState%22%3A%7B%22logs%22%3A%7B%22visualisationType%22%3A%22logs%22%7D%7D%7D%7D&orgId=1 to see the logs"
	@go build -o bin/backend cmd/backend/main.go && \
     	DEBUG=false bin/backend \
		| vector --config ./observe/vector.toml \
		&& echo -e "==> $(BLUE)Observe mode completed$(NC)" && make stop-observe \
		|| (echo -e "==> $(RED)Observe mode failed$(NC)" && make stop-observe && exit 1)

start-observe:
	@echo -e ":: $(GREEN)Starting observability stack...$(NC)"
	@docker compose -f observe/docker-compose.yml up --quiet-pull -d \
		&& echo -e "  -> Observability stack started" \
		|| (echo -e "  -> $(RED)Failed to start observability stack$(NC)" && exit 1)

stop-observe:
	$(call stop-observe)

build: gen
	@echo -e ":: $(GREEN)Building backend...$(NC)"
	@echo -e "  -> Building backend binary..."
	@go build -o bin/backend cmd/backend/main.go && echo -e "==> $(BLUE)Build completed successfully$(NC)" || (echo -e "==> $(RED)Build failed$(NC)" && exit 1)

gen:
	@echo -e ":: $(GREEN)Generating schema and code...$(NC)"
	@echo -e "  -> Running schema creation script..."
	@./scripts/create_full_schema.sh || (echo -e "  -> $(RED)Schema creation failed$(NC)" && exit 1)
	@echo -e "  -> Generating SQLC code..."
	@sqlc generate || (echo -e "  -> $(RED)SQLC generation failed$(NC)" && exit 1)
	@echo -e "  -> Running go generate..."
	@go generate ./... || (echo -e "  -> $(RED)Go generate failed$(NC)" && exit 1)
	@echo -e "==> $(BLUE)Generation completed$(NC)"

test: gen
	@echo -e ":: $(GREEN)Running tests...$(NC)"
	@go test -cover ./... && echo -e "==> $(BLUE)All tests passed$(NC)" || echo -e "==> $(RED)Tests failed$(NC)"

define wait_healthy
	@echo -e ":: $(GREEN)Waiting for $1 container to become healthy...$(NC)"
	@echo -e "  -> This may take a while..."
	@timeout 60 bash -c 'until [ "$$(docker inspect --format="{{.State.Health.Status}}" $1 2>/dev/null)" = "healthy" ]; do sleep 1; done' \
		&& echo -e "  -> $1 is healthy" \
		|| (echo -e "  -> $(RED)$1 did not become healthy in time$(NC)" && exit 1)
endef

define stop-observe
	@echo -e ":: $(GREEN)Stopping observability stack...$(NC)"
	@docker compose -f observe/docker-compose.yml down \
		&& echo -e "==> $(BLUE)Observability stack stopped$(NC)" \
		|| (echo -e "==> $(RED)Failed to stop observability stack$(NC)" && exit 1)
endef
