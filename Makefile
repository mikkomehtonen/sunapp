.PHONY: dev-backend dev-frontend test check clean

dev-backend:
	cd backend && go run cmd/server/main.go

dev-frontend:
	cd frontend && npm run dev

test:
	cd backend && go test ./internal/sun/... -v

check:
	@echo "=== Running backend tests ==="
	cd backend && go test ./internal/sun/... -v
	@echo ""
	@echo "=== Building backend ==="
	cd backend && go build ./cmd/server/
	@echo ""
	@echo "=== Checking frontend TypeScript ==="
	cd frontend && npx tsc --noEmit
	@echo ""
	@echo "=== Linting frontend ==="
	cd frontend && npx eslint .
	@echo ""
	@echo "All checks passed!"

clean:
	rm -f backend/server
	cd frontend && rm -rf dist
