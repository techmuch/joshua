ALTER TABLE claims ADD COLUMN archived BOOLEAN DEFAULT FALSE;

CREATE TABLE shares (
    id SERIAL PRIMARY KEY,
    solicitation_id INT REFERENCES solicitations(id) ON DELETE CASCADE,
    sender_id INT REFERENCES users(id),
    recipient_email TEXT NOT NULL,
    message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
