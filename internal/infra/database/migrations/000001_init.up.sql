 CREATE TABLE IF NOT EXISTS products
(
    id          UUID PRIMARY KEY,
    name        VARCHAR(100)   NOT NULL,
    description VARCHAR(500),
    price       DECIMAL(10, 2) NOT NULL,
    status      VARCHAR(20)    NOT NULL
);