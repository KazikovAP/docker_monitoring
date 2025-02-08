[![Go](https://img.shields.io/badge/-Go-464646?style=flat-square&logo=Go)](https://go.dev/)
[![TypeScript](https://img.shields.io/badge/-TypeScript-464646?style=flat-square&logo=TypeScript)](https://www.typescriptlang.org/)
[![React](https://img.shields.io/badge/-React-464646?style=flat-square&logo=React)](https://reactjs.org/)
[![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-464646?style=flat-square&logo=PostgreSQL)](https://www.postgresql.org/)
[![docker](https://img.shields.io/badge/-Docker-464646?style=flat-square&logo=docker)](https://www.docker.com/)
[![HTTP](https://img.shields.io/badge/-HTTP-464646?style=flat-square&logo=http)](https://developer.mozilla.org/en-US/docs/Web/HTTP)
[![Nginx](https://img.shields.io/badge/-Nginx-464646?style=flat-square&logo=Nginx)](https://www.nginx.com/)

# docker_monitoring
# Приложение для мониторинга состояния Docker-контейнеров

---
## Описание проекта
Проект представляет собой систему для постоянного мониторинга состояния Docker-контейнеров. Он состоит из нескольких компонентов, работающих вместе для отслеживания состояния контейнеров, пинга их IP-адресов и сохранения информации в базе данных для дальнейшего анализа и отображения на веб-странице.

---
## Компоненты системы
**1. Backend-сервис:**
    - Разработан с использованием языка Go.
    - Обеспечивает RESTful API для взаимодействия с фронтендом и базы данных.
    - API позволяет запросить данные о текущем состоянии контейнеров, а также добавлять новые записи о пингах в базу данных.

**2. Frontend-сервис:**
    - Разработан с использованием TypeScript и фреймворка React.
    - Получает данные о контейнерах через API Backend.
    - Отображает информацию о каждом контейнере в виде таблицы, включая IP-адрес, время последнего пинга и дату последней успешной попытки.
    - Для отображения данных используется библиотека для UI-компонентов Ant Design (antd).

**3. Сервис Pinger:**
    - Разработан на Go.
    - Получает список всех Docker-контейнеров.
    - Периодически пингует контейнеры и отправляет данные о результатах пинга в базу данных через API Backend.

**4. База данных PostgreSQL:**
    - Хранит данные о контейнерах, таких как IP-адреса, время пинга, статус контейнера, и дату последней успешной попытки пинга.
    - Обеспечивает доступ к данным через Backend-сервис.

---
## Технологии
* Go 1.23.0
* TypeScript/React
* PostgreSQL
* Docker
* REST API
* Nginx

---
## Запуск проекта

**1. Клонировать репозиторий:**
```
git clone https://github.com/KazikovAP/docker_monitoring.git
```

**2. Сборка и запуск проекта:**
```
docker-compose up --build
```

**3. Остановка и удаление контейнеров:**
```
docker-compose down
```

## Примеры запросов к API
### Получение списка пингов
**Пример GET запроса на адрес** `http://localhost:8080/pings`:

**Response:**
```JSON
[
  {
    "id": 1,
    "ip_address": "192.168.1.100",
    "ping_time": 120,
    "last_success_date": "2025-02-08T12:34:56Z",
    "created_at": "2025-02-08T12:34:56Z",
    "updated_at": "2025-02-08T12:34:56Z"
  },
  {
    "id": 2,
    "ip_address": "192.168.1.101",
    "ping_time": 140,
    "last_success_date": "2025-02-08T12:35:00Z",
    "created_at": "2025-02-08T12:35:00Z",
    "updated_at": "2025-02-08T12:35:00Z"
  }
]
```

### Добавление нового пинга
**Пример POST запроса на адрес с параметрами:**
```
http://localhost:8080/pings \
  -H "Content-Type: application/json" \
  -d '{
        "ip_address": "192.168.1.102",
        "ping_time": 150,
        "last_success_date": "2025-02-08T12:36:00Z"
      }'
```

**Response:**
```JSON
{
  "message": "Ping added successfully"
}
```

### Примеры ответов при ошибках
**Неверный формат IP-адреса (HTTP 400 Bad Request):**
**Response:**
```JSON
{
  "error": "Invalid IP address"
}
```

**Пинг с таким IP-адресом уже существует (HTTP 409 Conflict):**
**Response:**
```JSON
{
  "error": "IP address 192.168.1.102 already exists"
}
```

**Ошибка валидации JSON (HTTP 400 Bad Request):**
**Response:**
```JSON
{
  "error": "Key: 'Ping.IPAddress' Error:Field validation for 'IPAddress' failed on the 'required' tag"
}
```

---
## Разработал:
[Aleksey Kazikov](https://github.com/KazikovAP)

---
## Лицензия:
[MIT](https://opensource.org/licenses/MIT)
