
CREATE TABLE IF NOT EXISTS users (
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(100) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    role        VARCHAR(20) NOT NULL DEFAULT 'customer',
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);CREATE TABLE IF NOT EXISTS users (
    id          VARCHAR(36) PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    email       VARCHAR(100) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    role        VARCHAR(20) NOT NULL DEFAULT 'customer',
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Ligar customers a users
ALTER TABLE customers ADD COLUMN IF NOT EXISTS user_id VARCHAR(36) REFERENCES users(id);