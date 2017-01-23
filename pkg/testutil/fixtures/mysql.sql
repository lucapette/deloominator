DROP DATABASE IF EXISTS {{.Name}};
CREATE DATABASE {{.Name}};

USE {{.Name}};

CREATE TABLE event_types (
    id int NOT NULL,
    name varchar(255)
);

CREATE TABLE user_events (
    id int NOT NULL,
    name varchar(255),
    event_type_id int,
    user_id int
);

CREATE TABLE users (
    id int NOT NULL,
    name varchar(255)
);

INSERT INTO users (id, name) VALUES (42, 'Grace Hopper');
