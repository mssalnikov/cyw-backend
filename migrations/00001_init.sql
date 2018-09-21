-- +goose Up
CREATE TABLE users (
    id int NOT NULL PRIMARY KEY,
    fb_uid text,
    username text,
    email text
);

CREATE TABLE events (
    id int NOT NULL PRIMARY KEY,
    user_id int REFERENCES users(id) ON DELETE CASCADE,
    name text,
    description text,
    start_at timestamp without time zone,
    finish_at timestamp without time zone,
    is_private bool
);

CREATE TABLE points (
    id int NOT NULL PRIMARY KEY,
    event_id int references events(id) ON DELETE CASCADE,
    point geography(POINT),
    description text,
    have_question boolean,
    question text,
    answer text,
    is_chainded boolean,
    next_point int
);

-- +goose Down
DROP TABLE points;
DROP TABLE events;
DROP TABLE users;
