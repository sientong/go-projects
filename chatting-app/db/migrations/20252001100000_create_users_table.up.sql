CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,                     -- Unique identifier for each user
    username VARCHAR(100) NOT NULL UNIQUE,     -- Username, must be unique
    password VARCHAR(255) NOT NULL,            -- Password, hashed for security
    email VARCHAR(255) UNIQUE,                 -- Email address, optional but unique
    full_name VARCHAR(255),                    -- Full name of the user
    phone_number VARCHAR(20),                  -- Phone number
    is_active BOOLEAN DEFAULT TRUE,            -- Indicates if the user is active
    role VARCHAR(50) DEFAULT 'user',           -- User role (e.g., user, admin)
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
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Create index
CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_email ON users (email);