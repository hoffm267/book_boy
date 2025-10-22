-- OLD SCHEMA (without user_id) - for testing migrations
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- OLD SCHEMA: books without user_id
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    isbn TEXT UNIQUE NOT NULL,
    title TEXT,
    total_pages INTEGER
);

-- OLD SCHEMA: audiobooks without user_id
CREATE TABLE IF NOT EXISTS audiobooks (
    id SERIAL PRIMARY KEY,
    title TEXT,
    total_length INTERVAL
);

CREATE TABLE IF NOT EXISTS progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id INTEGER REFERENCES books(id) ON DELETE SET NULL,
    audiobook_id INTEGER REFERENCES audiobooks(id) ON DELETE SET NULL,
    book_page INTEGER,
    audiobook_time INTERVAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CHECK (
        book_id IS NOT NULL OR audiobook_id IS NOT NULL
    )
);

-- Test user (password is 'password123')
INSERT INTO users (username, email, password_hash) VALUES
  ('testuser', 'test@test.com', '$2a$10$Pk5vERjoku3o0IE4lnlvm.eNjryEDqUCJPYdU/VHOUKTtxA9HvFNO');

-- Old data without user_id
INSERT INTO books (isbn, title, total_pages) VALUES
  ('978-3-16-148410-0', 'Test Book', 300);

INSERT INTO audiobooks (title, total_length) VALUES
  ('Test Audiobook', INTERVAL '5 hours');
