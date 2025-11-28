## API Endpoints

### 1. Health Check
```
GET /api/health
```

### 2. Получить клиента по ID
```
GET /api/clients/{id}
```

### 3. Поиск клиентов
```
GET /api/clients/search?first_name=Иван&last_name=Иванов&birth_date=1990-01-01
```

Параметры (хотя бы один обязателен):
- `first_name` - имя клиента (частичное совпадение)
- `last_name` - фамилия клиента (частичное совпадение)
- `birth_date` - дата рождения (точное совпадение, формат YYYY-MM-DD)

### 4. Создать клиента
```
POST /api/clients
Content-Type: application/json

{
  "first_name": "Иван",
  "last_name": "Иванов",
  "birth_date": "1990-01-01",
  "core_data": {
    "phone": "+79991234567",
    "email": "ivan@example.com"
  },
  "features": {
    "income": 50000,
    "credit_history": "good"
  }
}
```

### 5. Обновить клиента
```
PUT /api/clients/{id}
Content-Type: application/json

{
  "first_name": "Иван",
  "last_name": "Иванов",
  "birth_date": "1990-01-01",
  "core_data": {...},
  "features": {...}
}
```

### 6. Удалить клиента
```
DELETE /api/clients/{id}
```

### 7. Рассчитать скоринг
```
GET /api/clients/{id}/scoring
```

Ответ:
```json
{
  "score": 0.85,
  "recommendations": [
    "Рекомендация 1",
    "Рекомендация 2"
  ],
  "factors": {
    "income": 0.3,
    "credit_history": 0.5,
    "age": 0.2
  }
}
```
