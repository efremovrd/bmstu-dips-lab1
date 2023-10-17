\c persons;

CREATE TABLE IF NOT EXISTS persons_ (
    id_ SERIAL PRIMARY KEY,
    name_ VARCHAR(64),
    address_ VARCHAR(64),
    work_ VARCHAR(64),
    age_ INT
);

GRANT ALL PRIVILEGES ON TABLE persons_ TO program;
GRANT ALL PRIVILEGES ON SEQUENCE persons__id__seq TO program;