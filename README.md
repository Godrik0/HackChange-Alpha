# HackChange-Alpha

Платформа управления клиентами с ML-скорингом.

## Быстрый старт

### Требования
- Docker
- Docker Compose

### Запуск

```bash
docker-compose up -d --build
```

### Доступ

- **Backend API**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **API Docs**: см. `backend/openapi.yml`

### Остановка

```bash
docker-compose down
```

### Удаление данных

```bash
docker-compose down -v
```