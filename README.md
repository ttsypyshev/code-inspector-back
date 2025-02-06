# Проверка кода студентов — Backend

Этот репозиторий содержит серверную часть проекта для проверки кода студентов. Он включает в себя API для взаимодействия с клиентской частью, работу с базой данных и интеграцию с другими сервисами.

## Репозитории
- **[Backend (вы здесь)](https://github.com/ttsypyshev/code-inspector-back)** – серверная часть проекта
- **[Frontend](https://github.com/ttsypyshev/code-inspector-front)** – клиентская часть проекта

## Ветки
- **[ssr](https://github.com/ttsypyshev/code-inspector-back/tree/ssr)** – разработка дизайна для 3 страниц
- **[database](https://github.com/ttsypyshev/code-inspector-back/tree/database)** – разработка структуры базы данных и ее подключение к бэкенду
- **[web-service](https://github.com/ttsypyshev/code-inspector-back/tree/web-service)** – создание веб-сервиса в бэкенде нашей системы для использования его в SPA
- **[auth](https://github.com/ttsypyshev/code-inspector-back/tree/auth)** – завершение бэкенда для SPA

## Установка и запуск

### 1. Клонирование репозитория
```bash
git clone https://github.com/ttsypyshev/code-inspector-back
cd code-inspector-back
```

### 2. Запуск приложения
```bash
sudo docker compose up -d
make run
```

### Ссылки для работы

- **Панель просмотра хранилища**  
  Для просмотра содержимого бакета `code-inspector` используйте следующую ссылку: [http://localhost:9001/buckets/code-inspector/browse](http://localhost:9001/buckets/code-inspector/browse)
  Перенаправление на хранилище: `http://localhost:9000`.

- **Браузер базы данных**  
  Для доступа к инструменту просмотра базы данных перейдите по ссылке (необходима настройка сервера и таблиц): [http://localhost:15432/browser/](http://localhost:15432/browser/)

- **Изучение методов API (Swagger)**  
  Для просмотра методов API, их request, response, curl, примеров используйте ссылку: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### Настройка базы данных

Для работы с базой данных используйте следующие SQL-скрипты:
- [Создать базу данных](doc/db/create.sql)
- [Вставить данные в базу данных](doc/db/insert.sql)
- [Удалить базу данных](doc/db/drop.sql)