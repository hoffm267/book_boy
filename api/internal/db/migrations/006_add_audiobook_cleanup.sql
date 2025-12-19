-- Migration: Add audiobook cleanup and improve book cleanup
-- Date: 2025-12-18

DROP TRIGGER IF EXISTS cleanup_orphaned_books ON progress;
DROP TRIGGER IF EXISTS cleanup_orphaned_books_on_delete ON progress;
DROP TRIGGER IF EXISTS cleanup_orphaned_audiobooks_on_delete ON progress;
DROP TRIGGER IF EXISTS cleanup_orphaned_books_on_update ON progress;
DROP TRIGGER IF EXISTS cleanup_orphaned_audiobooks_on_update ON progress;

CREATE OR REPLACE FUNCTION delete_orphaned_books()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.book_id IS NOT NULL THEN
        DELETE FROM books
        WHERE id = OLD.book_id
        AND NOT EXISTS (
            SELECT 1 FROM progress WHERE book_id = OLD.book_id
        );
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_orphaned_audiobooks()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.audiobook_id IS NOT NULL THEN
        DELETE FROM audiobooks
        WHERE id = OLD.audiobook_id
        AND NOT EXISTS (
            SELECT 1 FROM progress WHERE audiobook_id = OLD.audiobook_id
        );
    END IF;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER cleanup_orphaned_books_on_delete
AFTER DELETE ON progress
FOR EACH ROW
EXECUTE FUNCTION delete_orphaned_books();

CREATE TRIGGER cleanup_orphaned_audiobooks_on_delete
AFTER DELETE ON progress
FOR EACH ROW
EXECUTE FUNCTION delete_orphaned_audiobooks();

CREATE TRIGGER cleanup_orphaned_books_on_update
AFTER UPDATE ON progress
FOR EACH ROW
WHEN (OLD.book_id IS DISTINCT FROM NEW.book_id)
EXECUTE FUNCTION delete_orphaned_books();

CREATE TRIGGER cleanup_orphaned_audiobooks_on_update
AFTER UPDATE ON progress
FOR EACH ROW
WHEN (OLD.audiobook_id IS DISTINCT FROM NEW.audiobook_id)
EXECUTE FUNCTION delete_orphaned_audiobooks();
