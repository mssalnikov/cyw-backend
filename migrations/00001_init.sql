-- +goose Up
create database CYW ENCODING 'UTF-8' OWNER postgres;
grant all on database CYW to postgres;

CREATE TABLE users (
    id int NOT NULL PRIMARY KEY,
    fb_uid text,
    name text,
    last_name text
);

CREATE TABLE events (
    id int NOT NULL PRIMARY KEY,
    type text,
    user_id int references users(id)
);

CREATE TABLE points (
    id int NOT NULL PRIMARY KEY,
    event_id int references events(id),
    point geography(POINT),
    description text,
    have_question boolean,
    question text,
    answer text,
    is_chainded boolean,
    next_point int
);
-- INSERT INTO users VALUES
-- (0, 'root', '', ''),
-- (1, 'vojtechvitek', 'Vojtech', 'Vitek');

-- +goose Down
DROP TABLE points;
DROP TABLE events;
DROP TABLE users;
