CREATE TABLE params (
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    PRIMARY KEY(name)
);

INSERT INTO params (name, value) VALUES ('block_number', '0');

CREATE TABLE playerSaleOrders (
    id BIGINT NOT NULL,
    PRIMARY KEY(id)
);
