CREATE TYPE todo_status AS ENUM ('pending', 'completed');

CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    status todo_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL
);
