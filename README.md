# Organizational Structure API

REST API для управления организационной структурой компании.

## Стек

- Go
- net/http
- PostgreSQL
- GORM
- Goose
- Docker / Docker Compose

---

## Запуск проекта

```bash
docker compose up --build
```

## Миграции
Применить миграции:

```bash
goose -dir migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=organization_db sslmode=disable" up
```

## Тесты
Запуск тестов:
```bash
go test ./...
```
