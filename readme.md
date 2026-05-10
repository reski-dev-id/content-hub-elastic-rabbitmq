# Content Hub

Katalog Produk & Berita dengan full-text search.
Built with **Golang**, **MySQL**, **RabbitMQ**, dan **Elasticsearch** menggunakan **Clean Architecture** dan **Dependency Injection**.

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.22+ |
| HTTP Framework | Gin |
| Database | MySQL 8.0 + sqlx |
| Message Broker | RabbitMQ 3.x (amqp091-go) |
| Search Engine | Elasticsearch 8.x |
| DI Container | Wire (google/wire) |
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
- Wire mengelola seluruh dependency graph di `cmd/`

---

## Struktur Project

```
content-hub/
│
├── cmd/
│   ├── api/
│   │   ├── main.go           # entry point: wire inject + start server
│   │   ├── wire.go           # wire provider set
│   │   └── wire_gen.go       # wire generated (jangan diedit manual)
│   └── consumer/
│       ├── main.go           # entry point consumer goroutine
│       ├── wire.go
│       └── wire_gen.go
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
│   │       └── news.go
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
│   │       ├── middleware/
│   │       │   └── error.go
│   │       ├── request/      # request DTO + validator
│   │       │   ├── product.go
│   │       │   └── news.go
│   │       ├── response/     # response DTO
│   │       │   └── response.go
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
      "category":    { "type": "keyword" },
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
      "category":     { "type": "keyword" },
      "author":       { "type": "keyword" },
      "status":       { "type": "keyword" },
      "published_at": { "type": "date" }
    }
  }
}
```

---

## API Endpoints

| Method | Path | Query Params | Deskripsi |
|--------|------|--------------|-----------|
| `POST` | `/v1/products` | — | Insert produk |
| `GET` | `/v1/products` | `page`, `limit`, `category_id` | List + pagination dari MySQL |
| `GET` | `/v1/products/:id` | — | Detail produk |
| `PUT` | `/v1/products/:id` | — | Update produk |
| `POST` | `/v1/news` | — | Insert berita |
| `GET` | `/v1/news` | `page`, `limit`, `category_id` | List + pagination dari MySQL |
| `GET` | `/v1/news/:id` | — | Detail berita |
| `PUT` | `/v1/news/:id` | — | Update berita |
| `GET` | `/v1/search` | `q`, `type`, `category_id`, `page`, `limit` | Search via Elasticsearch |

### Contoh Request

```bash
# Insert product
curl -X POST http://localhost:8080/v1/products \
  -H "Content-Type: application/json" \
  -d '{"category_id":1,"title":"Laptop Gaming ASUS","price":15000000,"status":"active"}'

# List dengan pagination
curl "http://localhost:8080/v1/products?page=1&limit=10&category_id=1"

# Search by title + kategori
curl "http://localhost:8080/v1/search?q=laptop&type=product&category_id=1&page=1&limit=10"
```

---

## Setup & Menjalankan

### 1. Clone & konfigurasi

```bash
git clone https://github.com/reski-dev-id/content-hub-elastic-rabbitmq.git
cd content-hub
cp .env.example .env
```

### 2. Isi `.env`

```env
APP_PORT=8080
DB_HOST=localhost
DB_PORT=3306
DB_NAME=content_hub
DB_USER=root
DB_PASSWORD=secret
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
ELASTIC_URL=http://localhost:9200
```

### 3. Jalankan infrastruktur

```bash
docker-compose up -d
```

### 4. Install dependencies & generate Wire

```bash
go mod tidy
go install github.com/google/wire/cmd/wire@latest
cd cmd/api && wire && cd ../..
```

### 5. Migrasi database

```bash
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -path migrations -database "mysql://root:secret@tcp(localhost:3306)/content_hub" up
```

### 6. Buat Elasticsearch index

```bash
curl -X PUT http://localhost:9200/products \
  -H 'Content-Type: application/json' \
  -d @migrations/es_products.json

curl -X PUT http://localhost:9200/news \
  -H 'Content-Type: application/json' \
  -d @migrations/es_news.json
```


## ▶️ 7. Jalankan Aplikasi

Aplikasi ini membutuhkan **2 proses yang berjalan bersamaan**.

---

### Terminal 1 — API Server

go run cmd/api/main.go

Fungsi:
- menerima request (CRUD)
- menyimpan data ke MySQL
- menulis event ke `outbox_events`

---

### Terminal 2 — Outbox Consumer

go run cmd/consumer/main.go

Fungsi:
- membaca event dari `outbox_events`
- publish ke RabbitMQ
- update status event menjadi `sent`

---

## 🔁 Flow Sistem

Client Request  
↓  
API (Gin)  
↓  
MySQL (products + outbox_events)  
↓  
Outbox Poller (Consumer)  
↓  
RabbitMQ  

---

## ⚠️ Catatan Penting

- API dan Consumer **harus dijalankan bersamaan**
- Jika consumer tidak dijalankan:
  - data tetap masuk database
  - event tidak diproses
  - status tetap `pending`

---

## 🔍 Verifikasi

SELECT * FROM outbox_events;

Expected:
- sebelum consumer → pending
- setelah consumer → sent


### 8. Jalankan consumer (terminal terpisah)

```bash
go run cmd/consumer/main.go
```

---


## Prinsip Clean Architecture yang Diterapkan

| Prinsip | Implementasi |
|---------|-------------|
| Dependency Rule | Semua dependency arahnya ke dalam — domain tidak import infrastructure sama sekali |
| Interface Segregation | Tiap repository dan usecase punya interface tersendiri di `domain/` |
| Dependency Injection | Wire generate dependency graph otomatis dari constructor functions |
| Separation of Concern | Entity, business logic, data access, dan delivery sepenuhnya terpisah |
| Testability | Semua usecase dan handler mudah di-mock karena hanya bergantung pada interface |
| Single Responsibility | Tiap struct punya satu tanggung jawab yang jelas |
