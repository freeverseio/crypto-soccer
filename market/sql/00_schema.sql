CREATE TABLE playerSaleOrders (
    playerId BIGINT NOT NULL,
    PRIMARY KEY(playerId)
);

CREATE TABLE playerBuyOrders (
    playerId BIGINT NOT NULL REFERENCES playerSaleOrders(playerId),
    PRIMARY KEY(playerId)
)
