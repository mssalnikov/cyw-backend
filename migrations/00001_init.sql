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
    finish timestamp without time zone
);

CREATE TABLE points (
    id SERIAL NOT NULL PRIMARY KEY,
    event_id int references events(id) ON DELETE CASCADE,
    container text,
    naviaddress text,
    name text,
    question text,
    answer text,
    token int,
    prev_point_id int
);

CREATE TABLE userpoint (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id int references users(id) ON DELETE CASCADE,
    point_id int references points(id) ON DELETE CASCADE,
    is_found boolean,
    is_solved boolean
);

CREATE TABLE userevent (
    id SERIAL NOT NULL PRIMARY KEY,
    user_id int references users(id) ON DELETE CASCADE,
    event_id int references events(id) ON DELETE CASCADE,
    is_passed boolean
);

-- +goose Down
DROP TABLE userevent;
DROP TABLE userpoint;
DROP TABLE points;
DROP TABLE events;
DROP TABLE users;
