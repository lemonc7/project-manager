MIGRATE_PATH := ./database/migrate

# 加载.env文件
ifneq ("$(wildcard .env)","")
	include .env
	export
endif

DB_URL := postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable

ROOT_DB_URL := postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/postgres?sslmode=disable

db_create:
	@echo "🔍 Checking for database '${DB_NAME}'..."
	@if ! psql "${ROOT_DB_URL}" -tAc "SELECT 1 FROM pg_database WHERE datname = '${DB_NAME}'" | grep -q 1; then \
		echo "🆕 Database '${DB_NAME}' not found. Creating..."; \
		psql "${ROOT_DB_URL}" -c "CREATE DATABASE ${DB_NAME}"; \
	else \
		echo "✅ Database '${DB_NAME}' already exists."; \
	fi

db_drop:
	@echo "💥 Dropping database '${DB_NAME}'..."
	@if psql "${ROOT_DB_URL}" -tAc "SELECT 1 FROM pg_database WHERE datname = '${DB_NAME}'" | grep -q 1; then \
		psql "${ROOT_DB_URL}" -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname='${DB_NAME}';"; \
		psql "${ROOT_DB_URL}" -c "DROP DATABASE ${DB_NAME}"; \
		echo "✅ Database '${DB_NAME}' dropped."; \
	else \
		echo "ℹ️ Database '${DB_NAME}' does not exist."; \
	fi

# migrate_create name=xxx
migrate_create:
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: please use make migrate_create name=your_migration_name"; \
		exit 1; \
	fi
	@echo "📂 Creating migration: $(name)"
	migrate create -ext=sql -dir=$(MIGRATE_PATH) -seq $(name)
	@echo "✅ Migration created."

migrate_up:
	@echo "⬆️ Running migration UP..."
	migrate -path=${MIGRATE_PATH} -database=${DB_URL} up
	@echo "✅ Migration UP complete."

migrate_down:
	@echo "⬇️ Running migration DOWN (1 step)..."
	migrate -path=${MIGRATE_PATH} -database=${DB_URL} down 1
	@echo "✅ Migration DOWN complete."

migrate_drop:
	@echo "💥 Dropping all tables and the schema_migrations table..."
	migrate -path=${MIGRATE_PATH} -database=${DB_URL} drop

version:
	@echo "🏷️ Current migration version:"
	migrate -path=${MIGRATE_PATH} -database=${DB_URL} version