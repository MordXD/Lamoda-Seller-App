# Настройка Docker для проекта Kanban

## Структура проекта

Проект теперь полностью контейнеризован и включает:
- **Backend** (папка `back/`) - Go API сервер с Gin
- **Frontend** (папка `frontend/`) - React + Vite приложение
- **Database** - PostgreSQL
- **Migrations** - автоматические миграции БД

## Предварительные требования

1. Убедитесь, что у вас установлены Docker и Docker Compose
2. Создайте файл `.env` в корне проекта со следующими переменными:

```env
# Настройки базы данных
DB_USER=lamoda_user
DB_PASSWORD=lamoda_pass
DB_NAME=lamoda_db
DB_PORT=5431

# Настройки сервера
SERVER_PORT=8080

# JWT секрет для аутентификации
JWT_SECRET=supersecretkey
```

## Запуск проекта

1. **Сборка и запуск всех сервисов:**
```bash
docker-compose up --build
```

2. **Запуск в фоновом режиме:**
```bash
docker-compose up -d --build
```

3. **Остановка всех сервисов:**
```bash
docker-compose down
```

4. **Остановка с удалением данных:**
```bash
docker-compose down -v
```

## Доступ к приложению

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8080
- **Database:** localhost:5431

## API Эндпоинты

Backend предоставляет следующие эндпоинты:

### Публичные эндпоинты:
- `POST /auth/register` - Регистрация пользователя
- `POST /auth/login` - Вход в систему
- `GET /health` - Проверка здоровья сервиса

### Защищенные эндпоинты (требуют JWT токен):
- `GET /api/profile` - Получить профиль пользователя
- `PUT /api/profile` - Обновить профиль пользователя

## Что было настроено

### Backend (Go + Gin)
- Многоэтапная сборка Docker с оптимизацией
- Использование Alpine Linux для минимального размера
- Автоматические миграции GORM
- JWT аутентификация
- CORS настройки для фронтенда
- Health check эндпоинт

### Frontend (React + Vite)
- Многоэтапная сборка Docker
- Оптимизированный Nginx для SPA
- Gzip сжатие и кэширование
- Поддержка React Router

### База данных
- PostgreSQL 15 с проверкой здоровья
- Автоматическое создание пользователя и БД
- Персистентное хранение данных
- Порт 5431 (внешний доступ)

### Миграции
- Автоматический запуск при старте
- Ожидание готовности базы данных

## Разработка

Для разработки вы можете запускать сервисы по отдельности:

```bash
# Только база данных и миграции
docker-compose up db migrate

# Фронтенд в режиме разработки (в папке frontend)
cd frontend
yarn dev

# Бэкенд в режиме разработки (в папке back)
cd back
go run cmd/main.go
```

## Полезные команды

```bash
# Просмотр логов
docker-compose logs

# Просмотр логов конкретного сервиса
docker-compose logs app

# Пересборка конкретного сервиса
docker-compose build app

# Выполнение команды в контейнере
docker-compose exec app sh

# Просмотр таблиц в БД
docker-compose exec db psql -U lamoda_user -d lamoda_db -c "\dt"
```

## Тестирование API

Вы можете протестировать API с помощью curl или Postman:

```bash
# Проверка здоровья
curl http://localhost:8080/health

# Регистрация пользователя
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","password":"password123"}'

# Вход в систему
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
``` 