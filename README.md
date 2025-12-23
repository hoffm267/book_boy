# Book Boy: [bookboy.app](bookboy.app)

> A simple tracker that synchronizes progress between physical books and audiobooks

**Never lose your place when switching between reading and listening.**

## LOGIN FOR DEMO
Email: demo@bookboy.app
Password: Demo123!

## Summary

Book Boy automatically converts between book pages and audiobook timestamps. This enables users to switch between formats easily.

**Example**: You've read to page 150 of 300. Book Boy calculates you're at timestamp 4:59:00 in the 10 hour audiobook.

---

## Features

- **Bidirectional Sync**: Update book page → get audiobook time, or vice versa
- **JWT Authentication**: Secure user registration and login with 24 hour token expiration
- **ISBN Metadata Fetch**: Add books by ISBN, worker auto-fills title/pages via OpenLibrary API
- **RabbitMQ Workers**: Asynchronous background processing for metadata
- **Redis Caching**: Fast metadata lookups
- **REST API**: Endpoints for books, audiobooks, and progress tracking
- **Progress Filtering**: Track reading progress per user across multiple books

---

## Quick Start

### Prerequisites

- Docker >= 20.x (with `docker compose` support)
- Go 1.21+ (for local development)

### Running Locally

```bash
cd api

./scripts/dev.sh

make test

make clean
```

**Base URL**: `http://localhost:8080`
**Health Check**: `http://localhost:8080/health`

### Live API

**Production**: `http://3.146.159.19:8080`
**Frontend**: [Book Boy](https://book-boy-web.vercel.app)

### Deployment

Deployed on AWS EC2 with GitHub Actions for continuous deployment:
- Push to `main` → GitHub Actions builds Docker images
- Images pushed to Amazon ECR
- Automatic deployment to EC2 instance

---

## API Endpoints

### Authentication (Public)

```bash
POST /auth/register
Content-Type: application/json

{
  "username": "alice",
  "email": "alice@example.com",
  "password": "password123"
}

POST /auth/login
Content-Type: application/json

{
  "email": "alice@example.com",
  "password": "password123"
}

Response:
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": { ... }
}
```

### Protected Endpoints (Require JWT)

All endpoints require `Authorization: Bearer <token>` header.

**Books**
- `GET /books` - List all books
- `GET /books/:id` - Get specific book
- `POST /books` - Create new book (ISBN optional, auto-fills if provided)
- `PUT /books/:id` - Update book
- `DELETE /books/:id` - Delete book
- `GET /books/search?title=...` - Fuzzy search by title
- `GET /books/filter?...` - Filter books

**Audiobooks**
- `GET /audiobooks` - List all audiobooks
- `GET /audiobooks/:id` - Get specific audiobook
- `POST /audiobooks` - Create new audiobook
- `PUT /audiobooks/:id` - Update audiobook
- `DELETE /audiobooks/:id` - Delete audiobook
- `GET /audiobooks/search?title=...` - Fuzzy search by title

**Progress Tracking**
- `GET /progress` - List all progress
- `GET /progress/:id` - Get specific progress entry
- `POST /progress` - Create progress entry
- `PUT /progress/:id` - Update progress (auto-converts page ↔ time)
- `DELETE /progress/:id` - Delete progress
- `GET /progress/filter?...` - Filter by user/book/audiobook/status

**Metadata**
- `GET /metadata/isbn/:isbn` - Fetch book metadata by ISBN

**Quick Start Tracking**
- `POST /tracking/start` - Create book/audiobook and progress in one call
- `GET /tracking/current` - Get enriched progress with full book/audiobook details

**Users**
- `GET /users` - List all users
- `GET /users/:id` - Get specific user
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

**Health**
- `GET /health` - Check API health status

### Example: Add Book by ISBN

```bash
POST /books
Authorization: Bearer <token>

{
  "isbn": "9780316769174"
}

Response:
{
  "id": 1,
  "isbn": "9780316769174",
  "title": "",
  "total_pages": 0
}
```

Worker automatically fetches metadata and updates:
```bash
{
  "id": 1,
  "isbn": "9780316769174",
  "title": "The Catcher in the Rye",
  "total_pages": 277
}
```

### Example: Quick Start Tracking

```bash
POST /tracking/start
Authorization: Bearer <token>

{
  "format": "book",
  "title": "Atomic Habits",
  "author": "James Clear",
  "total_pages": 320,
  "current_page": 50
}

GET /tracking/current
Authorization: Bearer <token>

Response:
[
  {
    "progress_id": 1,
    "book": {
      "id": 1,
      "title": "Atomic Habits",
      "total_pages": 320
    },
    "current_page": 50,
    "completion_percent": 15,
    "updated_at": "2025-10-19T15:08:12Z"
  }
]
```

### Example: Manual Progress Tracking

```bash
POST /books
{
  "title": "The Great Gatsby",
  "total_pages": 180
}

POST /audiobooks
{
  "title": "The Great Gatsby",
  "total_length": "04:49:00"
}

POST /progress
{
  "book_id": 1,
  "audiobook_id": 1,
  "book_page": 90
}

Response:
{
  "id": 1,
  "book_id": 1,
  "audiobook_id": 1,
  "book_page": 90,
  "audiobook_time": "02:24:30"
}

PUT /progress/1
{
  "book_page": 135
}

Response:
{
  "book_page": 135,
  "audiobook_time": "03:36:45",
  "completion_percent": 75
}
```

---

## Testing

```bash
make test

go test -v ./internal/service -run TestAuthService

go test -cover ./...

./scripts/docker-test.sh
```

## Environment Variables

Create `.env` file (or set environment variables):

```bash
JWT_SECRET=your-secret-key-change-this-in-production
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSLMODE=disable
RABBITMQ_PASSWORD=your-rabbitmq-password
REDIS_URL=localhost:6379
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
BOOK_METADATA_SERVICE_URL=http://localhost:8000
```

---

## Architecture

**Monorepo Structure:**
```
book_boy/
├── api/                      # Go REST API
├── book_metadata_service/    # Python microservice (OpenLibrary API)
├── web/                      # React frontend
└── .github/workflows/        # CI/CD pipelines
```

**Infrastructure:**
- **API**: Go + Gin framework
- **Database**: PostgreSQL
- **Cache**: Redis
- **Queue**: RabbitMQ
- **Metadata**: Python FastAPI microservice
- **Deployment**: AWS EC2 + ECR
- **Frontend**: React + Vite (Vercel)
