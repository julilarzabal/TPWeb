# Variables
DB_CONTAINER_NAME := postgres-tpe
DB_USER := postgres

.PHONY: test pre-test run-tests post-test

# Target principal que se ejecuta con 'make test'
test:
	@make pre-test
	@make run-tests || (make post-test && exit 1)
	@make post-test

# 1. Tareas previas
pre-test:
	@sudo apt update && sudo apt install -y docker-compose-v2
	@echo "=== [1/3] Limpiando contenedores y volúmenes previos ==="
	@sudo docker compose down -v
	@echo "=== Generando código Go con sqlc ==="
	@sqlc generate
	@echo "=== Compilando el proyecto ==="
	@go build ./...
	@echo "=== Levantando contenedor de PostgreSQL ==="
	@sudo docker compose up -d
	@echo "=== Esperando a que PostgreSQL esté listo para recibir conexiones ==="
	@until sudo docker exec $(DB_CONTAINER_NAME) pg_isready -U $(DB_USER) > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "Base de datos lista."

# 2. Ejecución de tests
run-tests:
	@echo "=== [2/3] Ejecutando Tests ==="
	@go test -v ./...

# 3. Tareas posteriores
post-test:
	@echo "=== [3/3] Limpiando contenedores y volúmenes ==="
	@sudo docker compose down -v