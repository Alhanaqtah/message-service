DO $$ 
BEGIN
   IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'message_status') THEN
      CREATE TYPE message_status AS ENUM ('pending', 'processed');
   END IF;
END $$;

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content TEXT NOT NULL,
    status message_status DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP
);