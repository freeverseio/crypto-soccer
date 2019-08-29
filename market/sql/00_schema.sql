CREATE TABLE player_sale_orders (
    playerId BIGINT NOT NULL,
    price INT NOT NULL,
    PRIMARY KEY(playerId)
);

CREATE TABLE player_buy_orders (
    playerId BIGINT NOT NULL REFERENCES player_sale_orders(playerId),
    price INT NOT NULL,
    PRIMARY KEY(playerId)
)
