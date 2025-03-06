DROP TABLE IF EXISTS recipes CASCADE;
CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    cooking_time TEXT,
    description TEXT,
    instructions TEXT
);

DROP TABLE IF EXISTS ingredients; 
CREATE TABLE ingredients (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    recipe_id INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    unit TEXT NOT NULL,
    quantity REAL NOT NULL
);

DROP TABLE IF EXISTS recipe_tags;
CREATE TABLE recipe_tags (
    id SERIAL PRIMARY KEY,
    recipe_id INTEGER NOT NULL REFERENCES recipes (id) ON DELETE CASCADE,
    name TEXT NOT NULL
);
