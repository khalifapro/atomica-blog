-- ************************************** posts table
create table posts
(
    id         SERIAL PRIMARY KEY,
    title      TEXT NOT NULL,
    content    TEXT NOT NULL,
    tags       JSONB NULL,
    photos     JSONB NULL,
    updated_at DATE NOT NULL,
    created_at DATE NOT NULL,
    deleted_at DATE NULL
);

CREATE INDEX fkIdx_10 ON posts (title);

CREATE TABLE comments
(
    id  SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    creator varchar(50),
    post_id serial   NOT NULL,
    content text NOT NULL,
    updated_at date        NOT NULL,
    created_at date        NOT NULL,
    deleted_at date NULL,
    CONSTRAINT FK_11 FOREIGN KEY (post_id) REFERENCES posts ("id") ON DELETE cascade
);
