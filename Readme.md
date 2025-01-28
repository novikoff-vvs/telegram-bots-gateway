# Telegram Bot Gateway Service

## 📝 Описание
Данный сервис является входной точкой для Telegram бота, построенного на микросервисной архитектуре. Основной сервис выступает в роли шлюза (gateway), который принимает все обновления от Telegram Bot API и распределяет их по соответствующим микросервисам.

## 🏗 Архитектура

### 🔸 Основной сервис (Gateway)
- Обрабатывает все входящие обновления от Telegram Bot API
- Маршрутизирует запросы в соответствующие микросервисы
- Обеспечивает единую точку входа для всех взаимодействий с ботом

### 🔸 Микросервисы (/services)
Бизнес-логика распределена по отдельным микросервисам, каждый из которых отвечает за определенный функционал. Все микросервисы находятся в директории `/services`.

## 🔄 Коммуникация
Взаимодействие между основным сервисом и микросервисами реализовано через **gRPC**:
- ⚡️ Высокопроизводительный протокол обмена данными
- 📋 Строгая типизация благодаря Protocol Buffers
- 🚀 Эффективная сериализация/десериализация

## 📁 Структура проекта
```
.
├── app/main.go                    # Входная точка приложения
├── services/                  # Директория с микросервисами
│   ├── service1/             # Микросервис 1
│   │   └── internal/         # Внутренние зависимости
│   │       └── grpc/         # gRPC конфигурация
│   │   └── *.proto   # Proto файлы
│   │
│   ├── service2/             # Микросервис 2
│   │   └── internal/         # Внутренние зависимости
│   │       └── grpc/         # gRPC конфигурация
│   │   └── *.proto   # Proto файлы
│   │
│   └── ...                   # Другие микросервисы
```

## 🚀 Запуск проекта

### Предварительные требования
- Go 1.19 или выше
- Protocol Buffers
- Make (опционально)

### Установка и запуск
1. Клонируйте репозиторий
```bash
git clone <repository-url>
```

2. Перейдите в директорию проекта
```bash
cd <project-directory>
```

3. Запустите сервис
```bash
go run main.go
```

## 📝 Лицензия
MIT