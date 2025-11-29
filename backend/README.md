## API Endpoints

### Health Check
```
GET /health
```

### Clients

#### Создать клиента
```
POST /api/clients
```
Требуемые поля: `first_name`, `last_name`, `birth_date` (DD-MM-YYYY)  
Опциональные: `middle_name`, `features`

#### Получить клиента
```
GET /api/clients/{id}
```

#### Обновить клиента
```
PUT /api/clients/{id}
```
Все поля опциональны: `first_name`, `last_name`, `middle_name`, `birth_date`, `features`

#### Удалить клиента
```
DELETE /api/clients/{id}
```

#### Поиск клиентов
```
GET /api/clients/search?first_name={name}&last_name={surname}&birth_date={date}
```
Параметры (хотя бы один обязателен):
- `first_name` - частичное совпадение
- `last_name` - частичное совпадение  
- `birth_date` - точное совпадение (DD-MM-YYYY)

### Scoring

#### Рассчитать ML-скоринг
```
GET /api/clients/{id}/scoring
```
Возвращает: `score` (0-1), `recommendations`, `factors`

**Полная документация:** см. `openapi.yml`
