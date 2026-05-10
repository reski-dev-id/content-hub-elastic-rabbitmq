# Content Hub

Katalog Produk & Berita dengan full-text search.
Built with **Golang**, **MySQL**, **RabbitMQ**, dan **Elasticsearch** menggunakan **Clean Architecture**.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.25.1 |
| HTTP Framework | Gin |
| Database | MySQL 8.0 + sqlx |
| Message Broker | RabbitMQ 3.x (amqp091-go) |
| Search Engine | Elasticsearch 8.x |
| Dependency Injection | Manual Constructor Injection |
| Containerization | Docker + Docker Compose |

---

## Clean Architecture Layer

```
┌─────────────────────────────────────────────┐
│               Delivery Layer                │  ← HTTP Handler (Gin)
│           internal/delivery/http            │
├─────────────────────────────────────────────┤
│               Usecase Layer                 │  ← Business Logic
│             internal/usecase                │
├─────────────────────────────────────────────┤
│              Repository Layer               │  ← Data Access (interface)
│            internal/repository              │
├─────────────────────────────────────────────┤
│               Domain Layer                  │  ← Entity, Interface, DTO
│              internal/domain                │
└─────────────────────────────────────────────┘
         ↑ dependency hanya ke dalam ↑
```

**Aturan utama:**
- Tiap layer hanya boleh depend ke layer di bawahnya
- Domain layer tidak boleh import package luar (pure Go)
- Dependency ditanamkan via interface, bukan concrete struct

---

## Struktur Project

```
content-hub/
│
├── cmd/
│   ├── api/
│   │   └── main.go           # entry point API server
│   └── consumer/
│       └── main.go           # RabbitMQ consumer + Elasticsearch indexer
│
├── internal/
│   │
│   ├── domain/               # ★ Layer 1: Domain (pure, no external deps)
│   │   ├── entity/
│   │   │   ├── product.go    # Product struct
│   │   │   ├── news.go       # News struct
│   │   │   ├── category.go   # Category struct
│   │   │   └── outbox.go     # OutboxEvent struct
│   │   ├── repository/       # interface repository (contract)
│   │   │   ├── product.go
│   │   │   ├── news.go
│   │   │   ├── category.go
│   │   │   └── outbox.go
│   │   └── usecase/          # interface usecase (contract)
│   │       ├── product.go
│   │       ├── news.go
│   │       └── search.go
│   │
│   ├── repository/           # ★ Layer 2: Repository (implements domain/repository)
│   │   ├── mysql/
│   │   │   ├── product.go
│   │   │   ├── news.go
│   │   │   ├── category.go
│   │   │   └── outbox.go
│   │   └── elasticsearch/
│   │       ├── product.go
│   │       ├── news.go
│   │       └── search.go
│   │
│   ├── usecase/              # ★ Layer 3: Usecase (implements domain/usecase)
│   │   ├── product.go
│   │   ├── news.go
│   │   └── search.go
│   │
│   ├── delivery/             # ★ Layer 4: Delivery
│   │   └── http/
│   │       ├── handler/
│   │       │   ├── product.go
│   │       │   ├── news.go
│   │       │   └── search.go
│   │       └── router.go
│   │
│   ├── infrastructure/       # driver/adapter eksternal
│   │   ├── mysql/
│   │   │   └── db.go         # *sqlx.DB provider
│   │   ├── elasticsearch/
│   │   │   └── client.go     # *elasticsearch.Client provider
│   │   └── rabbitmq/
│   │       ├── connection.go # *amqp.Connection provider
│   │       ├── publisher.go
│   │       └── consumer.go
│   │
│   └── outbox/
│       └── poller.go         # goroutine outbox → RabbitMQ
│
├── migrations/
│   ├── 001_create_categories.sql
│   ├── 002_create_products.sql
│   ├── 003_create_news.sql
│   └── 004_create_outbox_events.sql
│
├── config/
│   └── config.go             # viper config loader
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── go.mod
└── README.md
```

---

## Database Schema (DDL)

```sql
CREATE TABLE categories (
  id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name       VARCHAR(100) NOT NULL,
  slug       VARCHAR(100) NOT NULL UNIQUE,
  type       ENUM('product', 'news') NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
  id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  category_id BIGINT UNSIGNED NOT NULL,
  title       VARCHAR(255) NOT NULL,
  slug        VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  price       DECIMAL(15,2) NOT NULL DEFAULT 0,
  status      ENUM('active', 'inactive') DEFAULT 'active',
  created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES categories(id),
  INDEX idx_title (title),
  INDEX idx_category (category_id)
);

CREATE TABLE news (
  id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  category_id  BIGINT UNSIGNED NOT NULL,
  title        VARCHAR(255) NOT NULL,
  slug         VARCHAR(255) NOT NULL UNIQUE,
  content      LONGTEXT,
  author       VARCHAR(100) NOT NULL,
  status       ENUM('draft', 'published') DEFAULT 'draft',
  published_at TIMESTAMP NULL,
  created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES categories(id),
  INDEX idx_title (title),
  INDEX idx_status_published (status, published_at)
);

CREATE TABLE outbox_events (
  id             BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  aggregate_type ENUM('product', 'news') NOT NULL,
  aggregate_id   BIGINT UNSIGNED NOT NULL,
  event_type     VARCHAR(50) NOT NULL,
  payload        JSON NOT NULL,
  status         ENUM('pending', 'sent') DEFAULT 'pending',
  created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  sent_at        TIMESTAMP NULL,
  INDEX idx_status (status)
);
```

---

## Elasticsearch Index Mapping

### Index: `products`

```json
{
  "mappings": {
    "properties": {
      "id":          { "type": "long" },
      "title":       { "type": "text", "analyzer": "standard" },
      "description": { "type": "text" },
      "category_id": { "type": "long" },
      "price":       { "type": "float" },
      "status":      { "type": "keyword" },
      "updated_at":  { "type": "date" }
    }
  }
}
```

### Index: `news`

```json
{
  "mappings": {
    "properties": {
      "id":           { "type": "long" },
      "title":        { "type": "text", "analyzer": "standard" },
      "content":      { "type": "text" },
      "category_id":  { "type": "long" },
      "author":       { "type": "keyword" },
      "status":       { "type": "keyword" },
      "published_at": { "type": "date" }
    }
  }
}
```

---

## API Endpoints

| Method | Path | Query Params | Description |
|--------|------|--------------|-----------|
| `POST` | `/v1/products` | — | Create product |
| `GET` | `/v1/products` | `page`, `limit`, `category_id` | List product |
| `GET` | `/v1/products/:id` | — | Product detail |
| `PUT` | `/v1/products/:id` | — | Update product |
| `DELETE` | `/v1/products/:id` | — | Delete product |
| `POST` | `/v1/news` | — | Create news |
| `GET` | `/v1/news` | `page`, `limit`, `category_id` | List news |
| `GET` | `/v1/news/:id` | — | News detail |
| `PUT` | `/v1/news/:id` | — | Update news |
| `DELETE` | `/v1/news/:id` | — | Delete news |
| `GET` | `/v1/search` | `q`, `type`, `category_id`, `page`, `limit` | Search via Elasticsearch |

---

## Setup & Menjalankan

### 1. Clone Repository

```bash
git clone https://github.com/reski-dev-id/content-hub-elastic-rabbitmq.git
cd content-hub
```

---

### 2. Build Docker

```bash
docker compose build
```

---

### 3. Jalankan Semua Service

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

## Flow Sistem

```
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
curl "http://localhost:8080/v1/search?q=iphone&type=product&page=1&limit=10"
```

---

## Prinsip Clean Architecture yang Diterapkan

| Prinsip | Implementasi |
|---------|-------------|
| Dependency Rule | Semua dependency arahnya ke dalam — domain tidak import infrastructure sama sekali |
| Interface Segregation | Tiap repository dan usecase punya interface tersendiri di `domain/` |
| Dependency Injection | Dependency diinject manual melalui constructor function |
| Separation of Concern | Entity, business logic, data access, dan delivery sepenuhnya terpisah |
| Testability | Semua usecase dan handler mudah di-mock karena hanya bergantung pada interface |
| Single Responsibility | Tiap struct punya satu tanggung jawab yang jelas |

