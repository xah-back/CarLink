# carLink

API для сервиса carLink, предназначенного для работы с пользователями, автомобилями и поездками.

## Описание

Backend API проекта carLink.
Система предоставляет серверную часть для управления пользователями, автомобилями, поездками, бронированиями и отзывами.

## Стек технологий

![Go](https://img.shields.io/badge/Go-1.25.1-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Web%20Framework-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-ORM-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-API%20Docs-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)

## Архитектура

```
                HTTP
┌────────┐ ─────────────> ┌────────────────────────┐
│ Client │                │   API Gateway (Gin)    │
└────────┘                │        routes          │
                           └───────────┬───────────┘
                                       │
        ┌──────────────────────────────┼────────────────────────────────┐
        │                              │                                │
        v                              v                                v
┌──────────────┐              ┌──────────────┐                 ┌──────────────┐
│ User Handler │              │  Car Handler │                 │ Trip Handler │
└──────┬───────┘              └──────┬───────┘                 └──────┬───────┘
       │                             │                                │
       v                             v                                v
┌──────────────┐              ┌──────────────┐                 ┌──────────────┐
│ User Service │              │ Car Service  │                 │ Trip Service │
└──────┬───────┘              └──────┬───────┘                 └──────┬───────┘
       │                             │                                │
       v                             v                                v
┌──────────────┐              ┌──────────────┐                 ┌──────────────┐
│ User Repo    │              │ Car Repo     │                 │ Trip Repo    │
└──────┬───────┘              └──────┬───────┘                 └──────┬───────┘
       │                             │                                │
       └──────────────┬──────────────┴───────────────┬────────────────┘
                      │                              │
                      v                              v
              ┌──────────────┐               ┌──────────────────┐
              │ Booking Repo │               │ Review Repo      │
              └──────┬───────┘               └─────────┬────────┘
                     │                                 │
                     v                                 v
              ┌──────────────┐               ┌──────────────────┐
              │BookingService│               │ Review Service   │
              └──────────────┘               └──────────────────┘

                           │
                           v
                   ┌──────────────────┐
                   │   PostgreSQL     │
                   │ (GORM +          │
                   │  AutoMigrate)    │
                   └──────────────────┘


        ┌─────────────────────────────────────────────┐
        │ Background Worker                           │
        │ trip_status_worker.go                       │
        │ (обновление статусов поездок)               │
        └─────────────────────────────────────────────┘

```

## Основные возможности

- Регистрация и авторизация пользователей
- Управление автомобилями
- Создание и управление поездками
- Бронирование поездок
- Отзывы и оценки
- Фоновое обновление статусов поездок
- Документация API через Swagger

### Запуск проекта

1. Клонирование репозитория

```bash
git clone https://github.com/dzakaev/CarLink.git
cd pharmacy
```

2. Запуск

```bash
make run
```

3. Если с MakeFile не запускается

```bash
go run cmd/app/main.go
```

## Перед запуском убедитесь, что:

- установлен Go (версия 1.20+)
- запущена база данных PostgreSQL
- в проекте настроены переменные окружения

## Разработчики

- [Усман (я)](https://github.com/dzakaev)
- [Зубайр](https://github.com/mutsaevz)
- [Али](https://github.com/var-go)
- [Асхаб](https://github.com/Askhab90)
