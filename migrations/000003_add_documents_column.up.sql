ALTER TABLE solicitations ADD COLUMN IF NOT EXISTS documents JSONB DEFAULT '[]'::jsonb;
