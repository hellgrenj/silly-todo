CREATE TABLE list (
    id serial PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE item (
    id serial PRIMARY KEY,
    name TEXT NOT NULL,
    done BOOLEAN NOT NULL DEFAULT false,
    list_id INTEGER REFERENCES list(id) ON DELETE CASCADE
);