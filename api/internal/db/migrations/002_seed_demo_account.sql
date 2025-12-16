-- Seed demo account for recruiters/demos
-- Credentials: demo@bookboy.app / Demo123!

INSERT INTO users (username, email, password_hash, created_at)
VALUES (
    'demo_user',
    'demo@bookboy.app',
    '$2a$10$GNYz4nFg7MlQrN5QhlLd0eZbn7Jp15LgBlZTgGPCqOvd4l8P2ju0O',
    CURRENT_TIMESTAMP
)
ON CONFLICT (email) DO NOTHING;

DO $$
DECLARE
    demo_user_id INTEGER;
BEGIN
    SELECT id INTO demo_user_id FROM users WHERE email = 'demo@bookboy.app';

    IF demo_user_id IS NULL THEN
        RAISE NOTICE 'Demo user already exists, skipping seed data';
        RETURN;
    END IF;

    INSERT INTO books (user_id, isbn, title, total_pages)
    VALUES
        (demo_user_id, '9780544003415', 'The Lord of the Rings', 1178),
        (demo_user_id, '9780451524935', '1984', 328),
        (demo_user_id, '9780061120084', 'To Kill a Mockingbird', 324),
        (demo_user_id, '9780143127550', 'Dune', 688),
        (demo_user_id, '9780316769174', 'The Catcher in the Rye', 277)
    ON CONFLICT (isbn) DO NOTHING;

    INSERT INTO audiobooks (user_id, title, total_length)
    VALUES
        (demo_user_id, 'Atomic Habits by James Clear', INTERVAL '5 hours 35 minutes'),
        (demo_user_id, 'The Hobbit by J.R.R. Tolkien', INTERVAL '11 hours 8 minutes'),
        (demo_user_id, 'Sapiens by Yuval Noah Harari', INTERVAL '15 hours 17 minutes')
    ON CONFLICT DO NOTHING;

    INSERT INTO progress (user_id, book_id, book_page, created_at, updated_at)
    SELECT
        demo_user_id,
        b.id,
        CASE
            WHEN b.title = 'The Lord of the Rings' THEN 450
            WHEN b.title = '1984' THEN 200
            WHEN b.title = 'Dune' THEN 150
        END,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    FROM books b
    WHERE b.user_id = demo_user_id
        AND b.title IN ('The Lord of the Rings', '1984', 'Dune')
    ON CONFLICT DO NOTHING;

    INSERT INTO progress (user_id, audiobook_id, audiobook_time, created_at, updated_at)
    SELECT
        demo_user_id,
        a.id,
        CASE
            WHEN a.title = 'Atomic Habits by James Clear' THEN INTERVAL '2 hours 30 minutes'
            WHEN a.title = 'The Hobbit by J.R.R. Tolkien' THEN INTERVAL '4 hours 15 minutes'
        END,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    FROM audiobooks a
    WHERE a.user_id = demo_user_id
        AND a.title IN ('Atomic Habits by James Clear', 'The Hobbit by J.R.R. Tolkien')
    ON CONFLICT DO NOTHING;

    RAISE NOTICE 'Demo account seeded successfully for user_id: %', demo_user_id;
END $$;
