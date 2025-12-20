-- Migration: Remove user_id from books and audiobooks
-- Date: 2025-12-19
-- Description: Books and audiobooks should be shared across users, only associated via progress

-- Drop indexes first
DROP INDEX IF EXISTS idx_books_user;
DROP INDEX IF EXISTS idx_audiobooks_user;

-- Remove user_id columns
ALTER TABLE books DROP COLUMN IF EXISTS user_id;
ALTER TABLE audiobooks DROP COLUMN IF EXISTS user_id;
