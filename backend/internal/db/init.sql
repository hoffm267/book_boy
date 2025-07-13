-- Create user table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL
);

-- Create books table
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    isbn TEXT UNIQUE NOT NULL,
    title TEXT
);

-- Create audiobooks table
CREATE TABLE IF NOT EXISTS audiobooks (
    id SERIAL PRIMARY KEY,
    isbn TEXT UNIQUE NOT NULL,
    title TEXT
);

-- Create user progress table (1:N)
CREATE TABLE IF NOT EXISTS user_book_progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id INTEGER REFERENCES books(id) ON DELETE SET NULL,
    audiobook_id INTEGER REFERENCES audiobooks(id) ON DELETE SET NULL,
    book_page INTEGER,
    audiobook_time INTERVAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Auto-update updated_at on row change
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = CURRENT_TIMESTAMP;
   RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON user_book_progress
FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

-- Indexes for performance
CREATE INDEX idx_user_progress_user ON user_book_progress(user_id);
CREATE INDEX idx_user_progress_book ON user_book_progress(book_id);
CREATE INDEX idx_user_progress_audio ON user_book_progress(audiobook_id);
