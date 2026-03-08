CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title TEXT,
    author TEXT,
    rating INT,
    year INT,
    description TEXT,
    type TEXT
);
