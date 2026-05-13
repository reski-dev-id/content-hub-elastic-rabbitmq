# Content Hub

Katalog Produk & Berita dengan full-text search.
Built with **Golang**, **MySQL**, **RabbitMQ**, dan **Elasticsearch** menggunakan **Clean Architecture**.

---

## Features

- Clean Architecture
- DTO Request Validation
- Swagger OpenAPI Documentation
- Full-text Search with Elasticsearch
- Fuzzy Search (Typo Tolerance)
- Recommendation Search System
- Outbox Pattern
- RabbitMQ Async Event Processing
- Graceful Shutdown
- Health Check Endpoint
- Pagination Meta Response
- Validation Error Formatter
- Structured Logging
- Request ID Middleware
- Recovery Middleware
- Dockerized Application
- GitHub Actions CI/CD
- DockerHub Auto Deploy

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.25.1 |
| HTTP Framework | Gin |
| Database | MySQL 8.0 + sqlx |
| Message Broker | RabbitMQ 3.x (amqp091-go) |
| Search Engine | Elasticsearch 8.x |
| API Documentation | Swagger / OpenAPI |
| Validation | go-playground/validator |
| Configuration | Viper |
| Logging | Zerolog |
| Dependency Injection | Manual Constructor Injection |
| Containerization | Docker + Docker Compose |
| CI/CD | GitHub Actions |
| Registry | DockerHub |

---

## Clean Architecture Layer

```txt
┌─────────────────────────────────────────────┐
│               Delivery Layer                │
│           internal/delivery/http            │
│                                             │
│  - Handler                                 │
│  - Middleware                              │
│  - Request DTO                             │
│  - Response Formatter                      │
└─────────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────┐
│               Usecase Layer                 │
│             internal/usecase                │
│                                             │
│  - Business Logic                           │
│  - Validation Flow                          │
│  - Recommendation Orchestration             │
└─────────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────┐
│              Repository Layer               │
│            internal/repository              │
│                                             │
│  - MySQL Repository                         │
│  - Elasticsearch Repository                 │
│  - RabbitMQ Publisher                       │
└─────────────────────────────────────────────┘
                     ↓
┌─────────────────────────────────────────────┐
│               Domain Layer                  │
│              internal/domain                │
│                                             │
│  - Entity                                   │
│  - Repository Interface                     │
│  - Usecase Interface                        │
└─────────────────────────────────────────────┘

         ↑ dependency hanya ke dalam ↑
```

### Aturan utama

- Tiap layer hanya boleh depend ke layer di bawahnya
- Domain layer tidak boleh import package luar (pure Go)
- Dependency ditanamkan via interface, bukan concrete struct
- Handler tidak langsung akses repository
- Usecase hanya bergantung pada interface repository

---

## Struktur Project

```txt
content-hub/
│
├── .github/
│   └── workflows/
│       └── deploy.yml
│
├── cmd/
│   ├── api/
│   │   └── main.go
│   │
│   └── consumer/
│       └── main.go
│
├── config/
│   └── config.go
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── internal/
│   │
│   ├── domain/
│   │   ├── entity/
│   │   │   ├── category.go
│   │   │   ├── product.go
│   │   │   ├── news.go
│   │   │   └── outbox.go
│   │   │
│   │   ├── repository/
│   │   │   ├── category.go
│   │   │   ├── product.go
│   │   │   ├── news.go
│   │   │   └── outbox.go
│   │   │
│   │   └── usecase/
│   │       ├── product.go
│   │       └── news.go
│   │
│   ├── delivery/
│   │   └── http/
│   │       ├── handler/
│   │       │   ├── health.go
│   │       │   ├── product.go
│   │       │   └── news.go
│   │       │
│   │       ├── middleware/
│   │       │   ├── logger.go
│   │       │   ├── recovery.go
│   │       │   └── request_id.go
│   │       │
│   │       ├── request/
│   │       │   ├── product.go
│   │       │   └── news.go
│   │       │
│   │       ├── response/
│   │       │   └── response.go
│   │       │
│   │       └── router.go
│   │
│   ├── helper/
│   │   ├── pagination.go
│   │   ├── parser.go
│   │   └── validation.go
│   │
│   ├── infrastructure/
│   │   ├── mysql/
│   │   │   └── db.go
│   │   │
│   │   ├── elasticsearch/
│   │   │   └── client.go
│   │   │
│   │   └── rabbitmq/
│   │       ├── connection.go
│   │       ├── publisher.go
│   │       └── consumer.go
│   │
│   ├── repository/
│   │   ├── mysql/
│   │   │   ├── category.go
│   │   │   ├── product.go
│   │   │   ├── news.go
│   │   │   └── outbox.go
│   │   │
│   │   └── elasticsearch/
│   │       ├── product.go
│   │       └── news.go
│   │
│   ├── usecase/
│   │   ├── product.go
│   │   └── news.go
│   │
│   ├── logger/
│   │   └── logger.go
│   │
│   └── outbox/
│       └── poller.go
│
├── migrations/
│   ├── 001_create_categories.sql
│   ├── 002_create_products.sql
│   ├── 003_create_news.sql
│   └── 004_create_outbox_events.sql
│
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── go.mod
└── README.md
```

---

## System Architecture

```txt
                    ┌──────────────┐
                    │    Client    │
                    └──────┬───────┘
                           │ HTTP
                           ▼
                 ┌──────────────────┐
                 │   Gin API Layer  │
                 └────────┬─────────┘
                          │
                          ▼
                 ┌──────────────────┐
                 │     Usecase      │
                 │  Business Logic  │
                 └────────┬─────────┘
                          │
            ┌─────────────┴─────────────┐
            ▼                           ▼
   ┌─────────────────┐         ┌─────────────────┐
   │      MySQL      │         │ Outbox Events   │
   └─────────────────┘         └────────┬────────┘
                                        │
                                        ▼
                              ┌─────────────────┐
                              │ Outbox Poller   │
                              └────────┬────────┘
                                       │
                                       ▼
                              ┌─────────────────┐
                              │    RabbitMQ     │
                              └────────┬────────┘
                                       │
                                       ▼
                              ┌─────────────────┐
                              │    Consumer     │
                              └────────┬────────┘
                                       │
                                       ▼
                              ┌─────────────────┐
                              │ Elasticsearch   │
                              └─────────────────┘
```

---

## Search Architecture

```txt
User Search Query
        │
        ▼
┌─────────────────────┐
│ Elasticsearch Query │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Multi Match Search  │
│ - title^3           │
│ - description       │
│ - content           │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Fuzziness AUTO      │
│ Typo Tolerance      │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Search Result       │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ Recommendation      │
│ more_like_this      │
└─────────────────────┘
```

---

## Search Features

### Full-text Search

Menggunakan Elasticsearch `multi_match` query.

### Fuzzy Search

Mendukung typo tolerance:

```txt
iphne   → iphone
androis → android
samsng  → samsung
```

Menggunakan:

```json
{
  "fuzziness": "AUTO"
}
```

### Recommendation Search

Menggunakan Elasticsearch:

```txt
more_like_this
```

Recommendation otomatis muncul berdasarkan hasil pertama pencarian.

---

## Database Schema (DDL)

```sql
CREATE TABLE categories (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) NOT NULL UNIQUE,
  type ENUM('product', 'news') NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  category_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  price DECIMAL(15,2) NOT NULL DEFAULT 0,
  status ENUM('active', 'inactive') DEFAULT 'active',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES categories(id),
  INDEX idx_title (title),
  INDEX idx_category (category_id)
);

CREATE TABLE news (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  category_id BIGINT UNSIGNED NOT NULL,
  title VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL UNIQUE,
  content LONGTEXT,
  author VARCHAR(100) NOT NULL,
  status ENUM('draft', 'published') DEFAULT 'draft',
  published_at TIMESTAMP NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES categories(id),
  INDEX idx_title (title),
  INDEX idx_status_published (status, published_at)
);

CREATE TABLE outbox_events (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  aggregate_type ENUM('product', 'news') NOT NULL,
  aggregate_id BIGINT UNSIGNED NOT NULL,
  event_type VARCHAR(50) NOT NULL,
  payload JSON NOT NULL,
  status ENUM('pending', 'sent') DEFAULT 'pending',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  sent_at TIMESTAMP NULL,
  INDEX idx_status (status)
);
```

---

## Elasticsearch Index Mapping

### Index: products

```json
{
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "title": { "type": "text", "analyzer": "standard" },
      "description": { "type": "text" },
      "category_id": { "type": "long" },
      "price": { "type": "float" },
      "status": { "type": "keyword" },
      "updated_at": { "type": "date" }
    }
  }
}
```

### Index: news

```json
{
  "mappings": {
    "properties": {
      "id": { "type": "long" },
      "title": { "type": "text", "analyzer": "standard" },
      "content": { "type": "text" },
      "category_id": { "type": "long" },
      "author": { "type": "keyword" },
      "status": { "type": "keyword" },
      "published_at": { "type": "date" }
    }
  }
}
```

---

## API Endpoints

| Method | Path | Query Params | Description |
|--------|------|--------------|-----------|
| POST | /v1/products | — | Create product |
| GET | /v1/products | page, limit, category_id | List product |
| GET | /v1/products/search | q, page, limit, category_id | Search product |
| GET | /v1/products/:id | — | Product detail |
| PUT | /v1/products/:id | — | Update product |
| DELETE | /v1/products/:id | — | Delete product |
| POST | /v1/news | — | Create news |
| GET | /v1/news | page, limit, category_id | List news |
| GET | /v1/news/search | q, page, limit, category_id | Search news |
| GET | /v1/news/:id | — | News detail |
| PUT | /v1/news/:id | — | Update news |
| DELETE | /v1/news/:id | — | Delete news |
| GET | /health | — | Health check |

---

## Swagger Documentation

Generate swagger docs:

```bash
swag init -g cmd/api/main.go
```

Open browser:

```txt
http://localhost:8080/v1/swagger/index.html
```

---

## Health Check

```bash
curl http://localhost:8080/health
```

Example response:

```json
{
  "success": true,
  "message": "service healthy",
  "data": {
    "status": "UP",
    "mysql": "UP",
    "elasticsearch": "UP",
    "rabbitmq": "UP"
  }
}
```

---

## Validation Error Format

```json
{
  "success": false,
  "message": "validation failed",
  "errors": {
    "title": "minimum length is 3",
    "price": "must be greater than 0"
  }
}
```

---

## Pagination Response Format

```json
{
  "success": true,
  "message": "products fetched successfully",
  "data": [],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

---

## Search Response Format

```json
{
  "success": true,
  "message": "products search fetched successfully",
  "data": {
    "items": [],
    "recommendations": []
  },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

---

## Setup & Menjalankan

### 1. Clone Repository

```bash
git clone https://github.com/reski-dev-id/content-hub-elastic-rabbitmq.git
cd content-hub-elastic-rabbitmq
```

---

### 2. Copy Environment

```bash
cp .env.example .env
```

---

### 3. Build Docker

```bash
docker compose build
```

---

### 4. Jalankan Semua Service

```bash
docker compose up
```

Service yang akan berjalan:

- MySQL
- RabbitMQ
- Elasticsearch
- API Server
- Outbox Poller
- RabbitMQ Consumer

---

## Docker Image

DockerHub:

```txt
programmerreski/content-hub
```

Pull latest image:

```bash
docker pull programmerreski/content-hub:latest
```

---

## CI/CD

Deployment flow:

```txt
feature/*
→ development
→ main
→ GitHub Actions
→ DockerHub auto push
```

Saat merge ke branch main:

- GitHub Actions otomatis build docker image
- Push image ke DockerHub
- Generate latest image tag

GitHub Actions workflow:

```txt
.github/workflows/deploy.yml
```

---

## Flow Sistem

```txt
Client Request
      ↓
API Server (Gin)
      ↓
MySQL + outbox_events
      ↓
Outbox Poller
      ↓
RabbitMQ
      ↓
Consumer
      ↓
Elasticsearch
```

---

## Verifikasi

### List Product

```bash
curl http://localhost:8080/v1/products
```

---

### Create Product

```bash
curl -X POST http://localhost:8080/v1/products \
-H "Content-Type: application/json" \
-d '{
  "category_id": 1,
  "title": "iPhone 16",
  "slug": "iphone-16",
  "description": "Apple smartphone",
  "price": 25000000,
  "status": "active"
}'
```

---

### Create News

```bash
curl -X POST http://localhost:8080/v1/news \
-H "Content-Type: application/json" \
-d '{
  "category_id": 6,
  "title": "AI Dunia",
  "slug": "ai-dunia",
  "content": "AI berkembang sangat cepat",
  "author": "Reski",
  "status": "published"
}'
```

---

### Search Product

```bash
curl "http://localhost:8080/v1/products/search?q=iphone&page=1&limit=10"
```

---

### Typo Search

```bash
curl "http://localhost:8080/v1/products/search?q=iphne"
```

---

## Prinsip Clean Architecture yang Diterapkan

| Prinsip | Implementasi |
|---------|-------------|
| Dependency Rule | Semua dependency arahnya ke dalam — domain tidak import infrastructure sama sekali |
| Interface Segregation | Tiap repository dan usecase punya interface tersendiri di domain |
| Dependency Injection | Dependency diinject manual melalui constructor function |
| Separation of Concern | Entity, business logic, data access, dan delivery sepenuhnya terpisah |
| Testability | Semua usecase dan handler mudah di-mock karena hanya bergantung pada interface |
| Single Responsibility | Tiap struct punya satu tanggung jawab yang jelas |
| Graceful Shutdown | Server menangani SIGTERM/SIGINT dengan proper cleanup |
| Async Processing | Outbox pattern + RabbitMQ untuk eventual consistency |

---

## Current Status

Project saat ini sudah memiliki:

- Production-ready backend foundation
- Async event-driven architecture
- Search indexing pipeline
- Standardized API response
- Docker-based deployment
- Automated CI/CD pipeline
- Health monitoring endpoint
- Swagger documentation
- Fuzzy typo-tolerant search
- Elasticsearch recommendation system
- Structured logging middleware
- Recovery middleware
- Request ID tracing
