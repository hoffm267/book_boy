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

CREATE INDEX idx_user_progress_user ON progress(user_id);
CREATE INDEX idx_user_progress_book ON progress(book_id);
CREATE INDEX idx_user_progress_audio ON progress(audiobook_id);
CREATE INDEX IF NOT EXISTS idx_books_title_trgm ON books USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_audiobooks_title_trgm ON audiobooks USING gin (title gin_trgm_ops);