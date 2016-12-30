SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

DROP DATABASE IF EXISTS {{.Name}};
CREATE DATABASE {{.Name}};

\c {{.Name}}

CREATE TABLE event_types (
    id integer NOT NULL,
    name character varying
);

CREATE TABLE user_events (
    id integer NOT NULL,
    name character varying,
    event_type_id integer,
    user_id integer
);


CREATE TABLE users (
    id integer NOT NULL,
    name character varying
);

ALTER TABLE ONLY event_types ADD CONSTRAINT event_types_pkey PRIMARY KEY (id);


ALTER TABLE ONLY user_events ADD CONSTRAINT user_events_pkey PRIMARY KEY (id);


ALTER TABLE ONLY users ADD CONSTRAINT users_pkey PRIMARY KEY (id);
