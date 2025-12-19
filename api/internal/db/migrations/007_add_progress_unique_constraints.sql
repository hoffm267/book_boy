-- Migration: Add unique constraints to prevent duplicate progress entries
-- Date: 2025-12-19

CREATE UNIQUE INDEX IF NOT EXISTS idx_progress_user_book
ON progress(user_id, book_id)
WHERE book_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_progress_user_audiobook
ON progress(user_id, audiobook_id)
WHERE audiobook_id IS NOT NULL;
