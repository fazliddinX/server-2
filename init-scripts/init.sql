CREATE EXTENSION postgres_fdw;

CREATE TABLE users
(
    id    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name  VARCHAR(50)  NOT NULL,
    email VARCHAR(100) NOT NULL,
    age   INT          NOT NULL
);

CREATE SERVER server2_fdw FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host '3.125.33.48', port '5433', dbname 'replica');
CREATE USER MAPPING FOR postgres SERVER server2_fdw OPTIONS (user 'postgres', password '123321');
CREATE FOREIGN TABLE users_server2 (
    id UUID,
    name VARCHAR(50),
    email VARCHAR(100),
    age INT
    ) SERVER server2_fdw OPTIONS (schema_name 'public', table_name 'users');