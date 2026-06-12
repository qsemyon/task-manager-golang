# Task manager

<img src="image.png" width="600">

Минималистичное веб-приложение для управления личными задачами с современным интерфейсом

## Технологический стек
* **Backend:** Go 1.26.4 с встроенным net/http
* **Frontend:** TypeScript + Vite
* **База данных:** PostgreSQL 18.4
* **Контейнеризация:** Docker

## Установка

```Bash
# Клонировать проект
git clone https://github.com/...
cd ...

# Создать .env файл
cat > .env << EOF
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=todo
DB_HOST=db
DB_PORT=5432
APP_PORT=8080
EOF

# Запустить приложение
docker-compose up --build
```
Приложение доступно на http://localhost:8080

## Функциональность
- Создание и удаление задач
- Отметить задачу как выполненную
- Сохранение данных в PostgreSQL
- Темная тема с красными акцентами
- Защита от XSS-атак (HTML escaping)

## API

| Метод | Endpoint | Описание | Тело запроса |
|-------|----------|---------|--------------|
| GET | `/tasks` | Получить все задачи | — |
| POST | `/tasks` | Создать новую задачу | `{"title": "Текст задачи"}` |
| PUT | `/tasks` | Обновить статус задачи | `{"id": 1, "done": true}` |
| DELETE | `/tasks?id=N` | Удалить задачу по ID | — |