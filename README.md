# simple-auth

JWT-аутентификация на Go с использованием Gin и PostgreSQL.

## Функционал

- Регистрация пользователей с хешированием паролей (bcrypt)
- Аутентификация и выдача JWT токенов
- Защищённые эндпоинты с middleware авторизации
- PostgreSQL для хранения данных
- Docker-ready для деплоя

## Стек технологий

- **Gin** - веб-фреймворк
- **PostgreSQL** - база данных
- **sqlx** - SQL toolkit
- **golang-jwt** - JWT токены авторизации
- **bcrypt** - хеширование паролей
- **envconfig** - конфигурация через переменные окружения
- **golang-migrate** - миграции БД

## Быстрый старт с Docker

```bash
git clone <repo-url>
cd simple-auth

cp .env.example .env

docker-compose up --build

migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/go-auth?sslmode=disable" up
```

Приложение доступно на `http://localhost:8000`

## Локальная разработка

### Требования

- Go 1.24+
- PostgreSQL 16+
- golang-migrate

### Установка

```bash
go mod download

brew install golang-migrate  # macOS
# или
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Настройка

Настроить `.env`:

```env
# Postgres
DB_HOST=localhost # или db, если запускаете в докер
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go-auth

# Server
SERVER_PORT=8000

# Auth
AUTH_SECRET_KEY=your-secret-key-here
AUTH_ACCESS_TOKEN_EXPIRES_MINUTES=60
AUTH_REFRESH_TOKEN_EXPIRES_DAYS=7
```

**Важно**: `ENVIRONMENT` не в `.env`! Задаётся вручную перед запуском.

### Миграции

```bash
createdb go-auth

migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/go-auth?sslmode=disable" up

migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/go-auth?sslmode=disable" down 1

migrate create -ext sql -dir migrations -seq migration_name
```

### Запуск

```bash
export ENVIRONMENT=LOCAL

go run cmd/main.go
```

## API Endpoints

### Публичные

**Регистрация пользователя**
```bash
POST /api/v1/users/
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}
```

**Аутентификация (получить JWT)**
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "testuser",
  "password": "password123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Защищённые (требуют JWT в заголовке)

**Получить секретные данные**
```bash
GET /api/v1/users/secret
Authorization: Bearer <jwt-token>
```

## Примеры запросов

```bash
# 1. Регистрация
curl -X POST http://localhost:8000/api/v1/users/ \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "test123"}'

# 2. Логин (получить токен)
TOKEN=$(curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "test123"}' \
  | jq -r '.token')

# 3. Доступ к защищённому эндпоинту
curl -X GET http://localhost:8000/api/v1/users/secret \
  -H "Authorization: Bearer $TOKEN"
```

## Структура проекта

```
simple-auth/
├── cmd/
│   └── main.go                 # Точка входа приложения
├── internal/
│   ├── config/                 # Загрузка конфигурации
│   ├── handlers/               # HTTP handlers
│   ├── middleware/             # Middleware (JWT auth)
│   ├── server/                 # Настройка сервера и роутинга
│   ├── service/                # Бизнес-логика
│   └── storage/
│       ├── repository/         # Слой доступа к данным
│       └── models/             # Модели БД
├── migrations/                 # SQL миграции
├── Dockerfile                  # Multi-stage Docker build
├── docker-compose.yml          # Оркестрация сервисов
├── .env.example                # Шаблон переменных окружения
└── README.md
```

## Docker команды

```bash
# Запустить все сервисы
docker-compose up -d

# Пересобрать и запустить
docker-compose up --build

# Остановить сервисы
docker-compose down

# Посмотреть логи
docker-compose logs -f backend

# Войти в контейнер БД
docker-compose exec db psql -U postgres -d go-auth

# Удалить volumes (БД будет очищена)
docker-compose down -v
```

## Переменные окружения

| Переменная | Описание | Пример |
|-----------|----------|--------|
| `ENVIRONMENT` | Окружение (LOCAL/PRODUCTION) | `LOCAL` |
| `DB_HOST` | Хост PostgreSQL | `localhost` |
| `DB_PORT` | Порт PostgreSQL | `5432` |
| `DB_USER` | Пользователь БД | `postgres` |
| `DB_PASSWORD` | Пароль БД | `postgres` |
| `DB_NAME` | Имя БД | `go-auth` |
| `SERVER_PORT` | Порт сервера | `8000` |
| `AUTH_SECRET_KEY` | Секретный ключ для JWT | `random-string` |
| `AUTH_ACCESS_TOKEN_EXPIRES_MINUTES` | Срок жизни токена (минуты) | `60` |

## Разработка

### Добавление новых миграций

```bash
# Создать файлы миграции
migrate create -ext sql -dir migrations -seq add_users_table

# Применить
migrate -path migrations -database "postgres://..." up
```
