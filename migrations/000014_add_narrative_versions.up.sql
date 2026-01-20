CREATE TABLE narrative_versions (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    content TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Backfill existing narratives as the first version
INSERT INTO narrative_versions (user_id, content, created_at)
SELECT id, narrative, NOW() FROM users WHERE narrative IS NOT NULL AND narrative != '';
