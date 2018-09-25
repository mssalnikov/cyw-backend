-- +goose Up
CREATE TABLE users (
    id SERIAL NOT NULL PRIMARY KEY,
    fb_uid text,
    username text,
    email text
);

CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    user_id int REFERENCES users(id) ON DELETE CASCADE,
    name text,
    description text,
    start timestamp without time zone,
    finish timestamp without time zone,
    is_private bool
);

CREATE TABLE points (
    id SERIAL NOT NULL PRIMARY KEY,
    event_id int references events(id) ON DELETE CASCADE,
    point geography(POINT),
    description text,
    naviaddress text,
    have_question boolean,
    question text,
    answer text,
    is_chainded boolean,
    next_point int
);

CREATE TABLE userpoint (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id int references users(id) ON DELETE CASCADE,
    point_id int references points(id) ON DELETE CASCADE,
    is_solved boolean
);

-- +goose Down
DROP TABLE userpoint;
DROP TABLE points;
DROP TABLE events;
DROP TABLE users;
