# 📚 Library API

Этот проект представляет собой REST API для управления книгами. Он включает следующие функции:

- 🔑 Авторизация с использованием Refresh и JWT токенов
- 🐘 Использование PosgreSQL для хранения информации
- ✉️ Подписка и отписка от рассылки
- 📬 Автоматическая отправка email-уведомлений
- 📦 Контейнеризация с Docker
- 🔄 Использование Kafka для событийного взаимодействия
- 📖 Swagger UI для удобной документации

## 🚀 Установка и запуск

### 1️⃣ Клонирование репозитория:
`git clone https://github.com/AgamariFF/Library.git`


### 2️⃣ Создание файла конфигурации `.env`
Создайте файл `.env` в `library/configs/` и добавьте в него:
```
SERVER_PORT=8080
DB_DSN=host=database user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Moscow
JWTCoo_expires_time_sec=1500
domain=localhost
jwtSecret=your_secret_key
SMTP_Name=example@inbox.ru
SMTP_Password=123456
```
<sub>Все значения указаны для примера<sub>

### 3️⃣ Запуск приложения в Docker

`docker-compose up --build`

### 4️⃣ Документация API (Swagger UI)

`http://localhost:8080/swagger/index.html`

## 📌 Основные эндпоинты API

### 🔹 Авторизация
- `POST /register` – Регистрация пользователя
- `POST /login` – Вход и получение refresh и JWT токенов
- `POST /logOut` – Выход из системы с удалением токенов

### 🔹 Управление книгами
- `GET /getBooks` – Получить список всех книг
- `GET /getBook` – Выдаёт всю информацию по переданному id книги в query параметрах (требуется аутентификация)
- `POST /addBook` – Добавить новую книгу (требуется аутентификация с правами администратора)
- `POST /modifyingBook` – Изменить данные уже существующей книги (требуется аутентификация с правами администратора)
- `DELETE /deleteBook` – Удалить книгу (требуется аутентификация с правами администратора)

### 🔹 Подписка на рассылку
- `POST /subscribe` – Подписаться на email-уведомления
- `POST /unsubscribe` – Отписаться от email-уведомлений

## 🛠 Технологии
- Golang + Gin
- PostgreSQL + GORM
- Kafka
- Docker + Docker Compose
- Swagger UI
- JWT для аутентификации
- SMTP для email рассылки