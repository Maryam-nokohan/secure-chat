-- Create a new database called 'DatabaseName'
-- Connect to the 'master' database to run this snippet
CREATE DATABASE user_db;

\c use_app

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    passhash TEXT NOT NULL,
    bio TEXT

);

