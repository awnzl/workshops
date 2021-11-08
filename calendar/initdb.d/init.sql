CREATE TABLE users(
    login    VARCHAR(64) NOT NULL PRIMARY KEY,
    timezone TEXT NOT NULL
);

CREATE TABLE events(
    id          VARCHAR(36) NOT NULL PRIMARY KEY,
    user_id     VARCHAR(64) REFERENCES users (login),
    title       TEXT NOT NULL,
    description TEXT,
    datetime    TIMESTAMP,
    duration    INT,
    notes       TEXT[]
);
