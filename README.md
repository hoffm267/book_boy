# Book Boy

> A simple tracker that synchronizes progress between physical books and audiobooks

**Never lose your place when switching between reading and listening.**


## Summary

Book Boy automatically converts between book pages and audiobook timestamps. This enables users to switch between formats easily.

**Example**: You've read to page 150 of 300. Book Boy calculates you're at timestamp 4:59:00 in the 10 hour audiobook.

---

## Features

- **Bidirectional Sync**: Update book page → get audiobook time, or vice versa
- **JWT Authentication**: Secure user registration and login with 24 hour token expiration
- **REST API**: Endpoints for books, audiobooks, and progress tracking
- **Similar Title Search**: Find books/audiobooks with fuzzy matching
- **Progress Filtering**: Track reading progress per user across multiple books
- **Test Coverage**: Comprehensive tests with Bruno API collections for manual testing

---

## Quick Start

### Prerequisites

- Docker >= 20.x (with `docker compose` support)
- Go 1.21+ (for local development)

### Running the Application

```bash
#MORE SCRIPTS BEING ADDED
# Start database and backend (from project root)
./scripts/book_boy

# Clean up (removes DB image, preserves volume)
make clean
```

**Base URL**: `http://localhost:8080`

---

## API Endpoints

### Authentication (Public)

```bash
# Register new user
POST /auth/register
Content-Type: application/json

{
  "username": "alice",
  "email": "alice@example.com",
  "password": "password123"
}

# Login (returns JWT token)
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
- `POST /books` - Create new book
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

**Quick Start Tracking**
- `POST /tracking/start` - Create book/audiobook and progress in one call
- `GET /tracking/current` - Get enriched progress with full book/audiobook details

**Users**
- `GET /users` - List all users
- `GET /users/:id` - Get specific user
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

### Example: Quick Start Tracking

```bash
# Start tracking a book (creates book + progress in one call)
POST /tracking/start
Authorization: Bearer <token>

{
  "format": "book",
  "title": "Atomic Habits",
  "author": "James Clear",
  "total_pages": 320,
  "current_page": 50
}

# View all current reading/listening
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
# Create book and audiobook separately
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

# Link them with progress tracking
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
  "audiobook_time": "02:24:30"  // Auto calculated
}

# Update progress (updates both formats)
PUT /progress/1
{
  "book_page": 135
}

Response:
{
  "book_page": 135,
  "audiobook_time": "03:36:45",  // Auto calculated
  "completion_percent": 75
}
```

---

### Run Unit Tests

```bash
# All tests
go test ./...

# Specific service (e.g., auth tests)
go test -v ./backend/internal/bl -run TestAuthService

# With coverage
go test -cover ./...
```

## Environment Variables

Create `.env` file (or set environment variables):

```bash
JWT_SECRET=your-secret-key-change-this-in-production
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=book_boy
```

---

## Roadmap

### Completed
- [x] JWT authentication with 24-hour token expiration
- [x] Cross-format progress sync (page ↔ timestamp conversion)
- [x] Fuzzy title search for books and audiobooks
- [x] Progress filtering by user/book/audiobook/status
- [x] Completion percentage calculation
- [x] Tracking endpoints for simplified workflow
- [x] Progress enrichment with book/audiobook details
- [x] Input validation (two-layer: binding + custom)
- [x] Custom error types

### Planned
- [ ] OpenLibrary API integration for auto populating book metadata
- [ ] Find/Create auddiobook information db
- [ ] Minimal web frontend
- [ ] Multi stage build for docker
- [ ] Pagination for large collections
- [ ] OpenAPI/Swagger documentation
- [ ] Reading statistics dashboard