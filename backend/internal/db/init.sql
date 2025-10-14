CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- TABLES
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    isbn TEXT UNIQUE NOT NULL,
    title TEXT,
    total_pages INTEGER
);

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

-- TRIGGERS
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON progress
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE OR REPLACE FUNCTION delete_orphaned_progress()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.book_id IS NULL AND NEW.audiobook_id IS NULL THEN
        DELETE FROM progress WHERE id = NEW.id;
        RETURN NULL;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER progress_cleanup_trigger
BEFORE UPDATE ON progress
FOR EACH ROW
EXECUTE FUNCTION delete_orphaned_progress();

-- INDEXES
CREATE INDEX idx_user_progress_user ON progress(user_id);
CREATE INDEX idx_user_progress_book ON progress(book_id);
CREATE INDEX idx_user_progress_audio ON progress(audiobook_id);
CREATE INDEX IF NOT EXISTS idx_books_title_trgm ON books USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_audiobooks_title_trgm ON audiobooks USING gin (title gin_trgm_ops);

-- DATA
-- USERS (password for all test users is 'password123')
INSERT INTO users (username, email, password_hash) VALUES
  ('alice', 'alice@example.com', '$2a$10$Pk5vERjoku3o0IE4lnlvm.eNjryEDqUCJPYdU/VHOUKTtxA9HvFNO'),
  ('bob', 'bob@example.com', '$2a$10$Pk5vERjoku3o0IE4lnlvm.eNjryEDqUCJPYdU/VHOUKTtxA9HvFNO'),
  ('carol', 'carol@example.com', '$2a$10$Pk5vERjoku3o0IE4lnlvm.eNjryEDqUCJPYdU/VHOUKTtxA9HvFNO');

-- BOOKS
INSERT INTO books (isbn, title, total_pages) VALUES
  ('978-3-16-148410-0', 'Go Programming Language', 400),
  ('978-0-13-110362-7', 'The C Programming Language', 274),
  ('978-0-201-03801-7', 'Design Patterns', 395);

-- AUDIOBOOKS
INSERT INTO audiobooks (title, total_length) VALUES
  ('Clean Code', INTERVAL '9 hours 30 minutes'),
  ('Refactoring', INTERVAL '7 hours 45 minutes'),
  ('Effective Java', INTERVAL '10 hours 15 minutes');

-- PROGRESS
INSERT INTO progress (user_id, book_id, audiobook_id, book_page, audiobook_time) VALUES
  (1, 1, NULL, 50, NULL),
  (2, NULL, 2, NULL, INTERVAL '1 hour 15 minutes'),
  (3, 3, NULL, 120, NULL);

