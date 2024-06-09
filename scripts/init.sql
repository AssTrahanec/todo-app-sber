CREATE TABLE todo_items (
                            id SERIAL PRIMARY KEY,
                            title VARCHAR(255) NOT NULL,
                            description TEXT,
                            date TIMESTAMP NOT NULL,
                            is_done BOOLEAN NOT NULL
);