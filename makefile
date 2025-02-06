.PHONY: run kill_process build frontend_build frontend_run backend_build backend_run debug


# Переменные
SRC_DIR = src
BACKEND_PORT = 8080
BACKEND_NAME = main
BACKEND_BUILD_DIR = bin
FRONTEND_DIR = frontend
FRONTEND_PORT = 3000
FRONTEND_NAME = node
FRONTEND_BUILD_DIR = dist

# Генерация Swagger-документации
swag:
	@echo "Генерация Swagger-документации..."
	@swag init --dir=$(SRC_DIR) --parseDependency --parseInternal --output ./docs/swagger
	@echo "Swagger-документация успешно сгенерирована."

# Убийство процесса, если нужно
kill_process:
	@echo "Поиск процесса '${BACKEND_NAME}' на порту ${BACKEND_PORT}..."
	@PID=$$(sudo lsof -t -i :${BACKEND_PORT} -sTCP:LISTEN | xargs ps -o pid=,comm= | grep ${BACKEND_NAME} | awk '{print $$1}') && \
	if [ -n "$$PID" ]; then \
		echo "Завершаем процесс '${BACKEND_NAME}' с PID: $$PID"; \
		sudo kill -9 $$PID; \
	else \
		echo "Процесс '${BACKEND_NAME}' не найден на порту ${BACKEND_PORT}."; \
	fi

	@echo "Поиск процесса '${FRONTEND_NAME}' на порту ${FRONTEND_PORT}..."
	@PID=$$(sudo lsof -t -i :${FRONTEND_PORT} -sTCP:LISTEN | xargs ps -o pid=,comm= | grep ${FRONTEND_NAME} | awk '{print $$1}') && \
	if [ -n "$$PID" ]; then \
		echo "Завершаем процесс с PID: $$PID"; \
		sudo kill -9 $$PID; \
	else \
		echo "Процесс '${FRONTEND_NAME}' не найден на порту ${FRONTEND_PORT}."; \
	fi


# Сборка Go-программы
backend_build: swag
	@echo "Сборка Go-программы..."
	@go build -v -o $(BACKEND_BUILD_DIR)/$(BACKEND_NAME) $(SRC_DIR)/main.go
	@echo "Сборка завершена. Программа сохранена в $(BACKEND_BUILD_DIR)/$(BACKEND_NAME)."

# Запуск Go-программы
# backend_run: kill_process backend_build
backend_run:
	@echo "Запуск Go-программы..."
#	@GIN_MODE=release ./$(EXECUTABLE)
	@go run $(SRC_DIR)/main.go

# Сборка фронтенда (React + TypeScript с Vite)
frontend_build:
	@echo "Сборка фронтенда..."
	@cd $(FRONTEND_DIR) && npm install && npm run build
	@echo "Сборка фронтенда завершена. Приложение собрано в $(FRONTEND_BUILD_DIR)."

# Запуск фронтенда (React + TypeScript с Vite)
frontend_run:
	@echo "Запуск фронтенда на порту ${FRONTEND_PORT}..."
	@cd $(FRONTEND_DIR) && npm run dev

# Общая сборка фронтенда и бэкенда
build: backend_build frontend_build
	@echo "Общая сборка завершена: бэкенд и фронтенд собраны."

# Запуск бэкенда и фронтенда вместе
run: kill_process
	@echo "Запуск бэкенда и фронтенда..."
	@$(MAKE) backend_run & 
	@$(MAKE) frontend_run &
	@wait


