CREATE DATABASE IF NOT EXISTS mydb;

USE mydb;

DROP TABLE IF EXISTS mydb.animals;
CREATE TABLE mydb.animals (
    id SERIAL PRIMARY KEY,
    name varchar(255),
    legs varchar(255)
);

CREATE USER 'read_user'@'%' IDENTIFIED BY '123456';
GRANT SELECT ON mydb.* TO 'read_user'@'%';

CREATE USER 'writeuser'@'%' IDENTIFIED BY 'passwordforwriteuser';
GRANT INSERT, UPDATE, SELECT, DELETE ON mydb.* TO 'writeuser'@'%';