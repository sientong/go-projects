CREATE TABLE IF NOT EXISTS channels (
    id SERIAL PRIMARY KEY,                     -- Unique identifier for each user
    channel_name VARCHAR(100) NOT NULL UNIQUE, -- Channel name, must be unique
    is_restricted BOOLEAN DEFAULT FALSE,       -- Indicates if the channel has limited access
    is_active BOOLEAN DEFAULT TRUE,            -- Indicates if the channels is active
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Record creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Last update timestamp
);

-- Create triggre for updated_at column using current timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_updated_at
BEFORE UPDATE ON channels
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();