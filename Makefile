.PHONY: dev-backend dev-frontend test check clean

.PHONY: build-frontend-dist
build-frontend-dist:
	cd frontend && npm run build
	rm -rf backend/internal/web/dist
	mkdir -p backend/internal/web
	cp -r frontend/dist backend/internal/web/dist

dev-backend:
	cd backend && go run cmd/server/main.go

dev-frontend:
	cd frontend && npm run dev -- --host 0.0.0.0

test:
	cd backend && go test ./internal/sun/... -v

check: build-frontend-dist
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
	rm -rf backend/internal/web/dist
	rm -rf frontend/dist
