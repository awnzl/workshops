CREATE TABLE users(
    id       VARCHAR(32) NOT NULL PRIMARY KEY,
    login    VARCHAR(64) NOT NULL,
    timezone TEXT NOT NULL
);

CREATE TABLE events(
    id          VARCHAR(32) NOT NULL PRIMARY KEY,
    user_id     VARCHAR(32) REFERENCES users (id),
    title       TEXT NOT NULL,
    description TEXT,
    time        TIMESTAMP NOT NULL,
    timezone    TEXT NOT NULL,
    duration    INT,
    notes       TEXT[]
);
