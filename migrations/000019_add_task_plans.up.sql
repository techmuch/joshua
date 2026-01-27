ALTER TABLE tasks ADD COLUMN plan TEXT;
ALTER TABLE tasks ADD COLUMN plan_status TEXT DEFAULT 'none';

CREATE TABLE task_comments (
    id SERIAL PRIMARY KEY,
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
