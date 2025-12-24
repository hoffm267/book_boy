# Book Boy

> Cross-format reading tracker that synchronizes progress between physical books and audiobooks

**Live Demo**: [bookboy.app](https://bookboy.app)

Never lose your place when switching between reading and listening.

---

## What It Does

Book Boy automatically converts between book pages and audiobook timestamps, enabling seamless format switching.

**Example**: You've read to page 150 of a 300-page book. Book Boy calculates you're at 4:59:00 in the 10-hour audiobook. Switch to your commute, start listening at 4:59:00, then pick up your physical book later at the exact page where you left off.

### Key Features

- **Bidirectional Sync**: Update book page → get audiobook time, or vice versa
- **Real-time Updates**: SSE (Server-Sent Events) push metadata changes to frontend instantly
- **Async Metadata Enrichment**: RabbitMQ workers fetch book data from external APIs in background
- **Redis Caching**: Fast metadata lookups with 10-minute TTL
- **JWT Authentication**: Secure user sessions with 24-hour token expiration
- **Shared Resources**: Books/audiobooks are shared across users, progress is user-specific

---

## Quick Start

### Prerequisites

- Docker >= 20.x with `docker compose` support

### Run Locally

```bash
# Clone the repository
git clone https://github.com/yourusername/book_boy.git
cd book_boy

# Start all services (API, web, database, Redis, RabbitMQ, metadata service)
docker compose up -d

# View logs
docker compose logs -f api
docker compose logs -f web

# Stop all services
docker compose down
```

**Access Points:**
- **Frontend**: http://localhost:5173
- **API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **RabbitMQ Management**: http://localhost:15672 (guest/guest)

**Demo Login:**
- Email: `demo@bookboy.app`
- Password: `Demo123!`

### Run Tests

```bash
cd api
go test ./...
```

---

### Architecture

**Microservices:**
- **API** (Go): REST API with JWT auth, progress tracking, cross-format sync
- **Web** (React): SPA frontend with real-time SSE updates
- **Metadata Service** (Python): Async book metadata fetching from external APIs
- **Database** (PostgreSQL): User data, books, audiobooks, progress
- **Cache** (Redis): 10-minute TTL for book metadata
- **Queue** (RabbitMQ): Event-driven metadata enrichment

**Data Model:**
- Books/audiobooks are **shared resources** (no user_id)
- Progress records are **user-specific** (can link to both book_id and audiobook_id)
- Single progress record enables cross-format sync

**Event-Driven Flow:**
1. User creates book → API publishes event to RabbitMQ
2. Python service consumes event → Fetches metadata from external APIs
3. Service publishes result → API caches in Redis
4. API broadcasts SSE event → Frontend updates in real-time

---

## API Endpoints

### Authentication

```bash
POST /auth/register
POST /auth/login
```

### Endpoints

All endpoints require `Authorization: Bearer <token>` header.

**Books**
- `GET /books` - List all books
- `POST /books` - Create book (ISBN auto-fills metadata via worker)
- `PUT /books/:id` - Update book
- `DELETE /books/:id` - Delete book
- `GET /books/search?title=...` - Fuzzy search

**Audiobooks**
- `GET /audiobooks` - List all audiobooks
- `POST /audiobooks` - Create audiobook
- `PUT /audiobooks/:id` - Update audiobook
- `DELETE /audiobooks/:id` - Delete audiobook
- `GET /audiobooks/search?title=...` - Fuzzy search

**Progress Tracking**
- `GET /progress` - List user's progress
- `GET /progress/enriched` - List with full book/audiobook data (single query)
- `POST /progress` - Create progress entry
- `PUT /progress/:id` - Update progress (auto-converts page ↔ time)
- `DELETE /progress/:id` - Delete progress

**Tracking**
- `POST /tracking/start` - Create book/audiobook + progress in one call
- `GET /tracking/current` - Get current reading list with enriched data

**Real-time**
- `GET /events?token=<jwt>` - SSE stream for metadata updates

**Health**
- `GET /health` - API health check

### Example: Track Progress

```bash
# Login
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@bookboy.app","password":"Demo123!"}'

# Response: {"token":"eyJ...","user":{...}}

# Create book
curl -X POST http://localhost:8080/books \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"The Great Gatsby","author":"F. Scott Fitzgerald","total_pages":180,"isbn":"9780743273565"}'

# Create audiobook
curl -X POST http://localhost:8080/audiobooks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"title":"The Great Gatsby","author":"F. Scott Fitzgerald","narrator":"Jake Gyllenhaal","total_length":"04:49:00"}'

# Create progress (links both book and audiobook)
curl -X POST http://localhost:8080/progress \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"book_id":1,"audiobook_id":1,"book_page":90}'

# Response: {"id":1,"book_id":1,"audiobook_id":1,"book_page":90,"audiobook_time":"02:24:30"}

# Update progress - audiobook time auto-calculated
curl -X PUT http://localhost:8080/progress/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"book_page":135}'

# Response: {"book_page":135,"audiobook_time":"03:36:45","completion_percent":75}
```

---

## Environment Variables

### Local Development (docker-compose.yml)

**No `.env` file needed!** docker-compose has sensible defaults for local development.

Optional: Create `.env` to override defaults:
```bash
JWT_SECRET=my-custom-secret
DB_PASSWORD=my-custom-password
```

---

## Deployment

**Production**: Deployed on AWS EC2 with GitHub Actions CI/CD

**Production URLs:**
- **Frontend**: https://bookboy.app
- **API**: https://api.bookboy.app

---

## Tech Stack

- **API**: Go 1.24 + Gin framework
- **Frontend**: React + Vite + Bun
- **Database**: PostgreSQL 14
- **Cache**: Redis 7
- **Queue**: RabbitMQ 3
- **Metadata Service**: Python 3.12
- **Web Server**: Nginx (production)
- **Deployment**: Docker multi-stage builds, AWS EC2, GitHub Actions
- **Testing**: Go test, Bruno API collections

---

## Project Structure

```
book_boy/
├── api/                        # Go REST API
│   ├── cmd/server/            # Entry point
│   ├── internal/              # Services, repos, controllers
│   └── bruno/                 # API test collections
├── web/                       # React frontend
│   └── src/
├── book_metadata_service/     # Python microservice
├── .github/workflows/         # CI/CD
├── docker-compose.yml
└── docker-compose.production.yml
```

---

## Testing

### API Tests
```bash
cd api
go test ./...
go test -v ./internal/service
go test -cover ./...
```

### API Integration Tests (Bruno)
1. Install [Bruno](https://www.usebruno.com/)
2. Open Bruno and load collections from `/api/bruno/`
3. Use the `local` environment
4. Run collections for authentication, books, audiobooks, progress

---

## License

MIT
