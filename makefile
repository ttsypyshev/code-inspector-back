.PHONY: run kill_process build

# Переменные
PORT = 8080
PROCESS_NAME = main
BUILD_DIR = bin
SRC_DIR = src
EXECUTABLE = $(BUILD_DIR)/$(PROCESS_NAME)

# Генерация Swagger-документации
swag:
	@echo "Генерация Swagger-документации..."
	@swag init --dir=$(SRC_DIR) --parseDependency --parseInternal --output ./docs/swagger
	@echo "Swagger-документация успешно сгенерирована."

# Убийство процесса, если нужно
kill_process:
	@echo "Поиск процесса '${PROCESS_NAME}' на порту ${PORT}..."
	@PID=$$(sudo lsof -t -i :${PORT} -sTCP:LISTEN | xargs ps -o pid=,comm= | grep ${PROCESS_NAME} | awk '{print $$1}') && \
	if [ -n "$$PID" ]; then \
		echo "Завершаем процесс '${PROCESS_NAME}' с PID: $$PID"; \
		sudo kill -9 $$PID; \
	else \
		echo "Процесс '${PROCESS_NAME}' не найден на порту ${PORT}."; \
	fi

# Запуск в Docker
docker-build:
	@echo "Сборка Docker-образа..."
	@sudo docker compose build backend

docker-run: docker-build
	@echo "Запуск Go-программы в Docker..."
	@sudo docker compose up backend

docker-stop:
	@echo "Остановка контейнера..."
	@sudo docker compose down

# Сборка Go-программы
build: swag
	@echo "Сборка Go-программы..."
	@go build -v -o $(EXECUTABLE) $(SRC_DIR)/main.go
	@echo "Сборка завершена. Программа сохранена в $(EXECUTABLE)."

# Запуск Go-программы
run: kill_process build
	@echo "Запуск Go-программы..."
	@GIN_MODE=release ./$(EXECUTABLE)

# Запуск Go-программы
debug: kill_process
	@echo "Запуск Go-программы (в debug режиме)..."
	@go run $(SRC_DIR)/main.go