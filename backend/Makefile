GREEN = \033[0;32m
BLUE = \033[0;34m
RED = \033[0;31m
NC = \033[0m

all: gen build test
	@echo -e "==> $(BLUE)All tasks completed successfully$(NC)"

run: gen
	@echo -e ":: $(GREEN)Starting backend...$(NC)"
	@go run cmd/backend/main.go

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
