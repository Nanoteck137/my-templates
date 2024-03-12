-- +goose Up
CREATE TABLE todos (
    id TEXT,
    content TEXT NOT NULL,
    done BOOL NOT NULL,

    CONSTRAINT todos_pk PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE todos;
