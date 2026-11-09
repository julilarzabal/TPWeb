# Variables
DB_CONTAINER_NAME := postgres-tpe

.PHONY: run pre-run run-app post-run

# Target principal que se ejecuta con 'make run'
run:
	@make pre-run
	@make run-app || (make post-run && exit 1)
	@make post-run

# 1. Tareas previas
pre-run:
	@echo "=== [1/3] Limpiando contenedores y volúmenes previos ==="
	@sudo docker compose down -v
	@kill $$(cat /tmp/tpe-server.pid 2>/dev/null) 2>/dev/null || true
	@rm -f /tmp/tpe-server.pid
	@echo "=== Generando código Go con sqlc ==="
	@sqlc generate
	@echo "=== Generando código Go con templ ==="
	@go run github.com/a-h/templ/cmd/templ@v0.3.1020 generate
	@echo "=== Compilando el proyecto ==="
	@go build -o /tmp/tpe-server .
	@echo "=== Levantando contenedor de PostgreSQL ==="
	@sudo docker compose up -d
	@echo "=== Esperando a que PostgreSQL esté listo para recibir conexiones ==="
	@until sudo docker exec $(DB_CONTAINER_NAME) sh -c 'psql -U "$$POSTGRES_USER" -d "$$POSTGRES_DB" -c "SELECT pg_sleep(2)"' > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "Base de datos lista."

# 2. Ejecución de la aplicación
run-app:
	@echo "=== [2/3]  Ejecutando la aplicación ==="
	@/tmp/tpe-server > /tmp/tpe-server.log 2>&1 & echo $$! > /tmp/tpe-server.pid
	@echo "=== Esperando a que el servidor Go esté listo ==="
	@until nc -z localhost 8080 > /dev/null 2>&1; do \
		kill -0 $$(cat /tmp/tpe-server.pid) 2>/dev/null || { echo "El servidor Go terminó antes de estar listo:"; cat /tmp/tpe-server.log; exit 1; }; \
		sleep 1; \
	done
	@echo "Servidor Go listo en http://localhost:8080/."
	@echo "=== Abriendo la aplicación en el navegador ==="
	@xdg-open http://localhost:8080/ > /dev/null 2>&1 || open http://localhost:8080/ > /dev/null 2>&1 || echo "Abrí manualmente http://localhost:8080/."
	@echo "Presioná Enter para detener la aplicación y limpiar los contenedores..."
	@read -r _ || true

# 3. Tareas posteriores
post-run:
	@echo "=== [3/3] Limpiando contenedores y volúmenes ==="
	@sudo docker compose down -v
	@echo "=== Apagando servidor Go ==="
	@kill $$(cat /tmp/tpe-server.pid 2>/dev/null) 2>/dev/null || true
	@rm -f /tmp/tpe-server.pid /tmp/tpe-server