CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    cards JSONB DEFAULT '[]',
    map_info JSONB DEFAULT '{"point":{"lat":55.753544,"lon":37.621202},"radius":3000}',
    exchange INTEGER DEFAULT 600
);
