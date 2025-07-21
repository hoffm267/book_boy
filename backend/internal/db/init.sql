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

-- Sample data
INSERT INTO users (username) VALUES
  ('alice'),
  ('bob'),
  ('carol');

INSERT INTO books (isbn, title) VALUES
  ('978-3-16-148410-0', 'Go Programming Language'),
  ('978-0-13-110362-7', 'The C Programming Language'),
  ('978-0-201-03801-7', 'Design Patterns');

INSERT INTO audiobooks (isbn, title) VALUES
  ('978-1-60309-452-8', 'Clean Code'),
  ('978-0-596-52068-7', 'Refactoring'),
  ('978-0-321-63537-8', 'Effective Java');

INSERT INTO user_book_progress (user_id, book_id, audiobook_id, book_page, audiobook_time) VALUES
  (1, 1, NULL, 50, NULL),
  (2, NULL, 2, NULL, INTERVAL '1 hour 15 minutes'),
  (3, 3, NULL, 120, NULL);
