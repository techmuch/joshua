CREATE TABLE solicitation_comments (
    id SERIAL PRIMARY KEY,
    solicitation_id INT REFERENCES solicitations(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
