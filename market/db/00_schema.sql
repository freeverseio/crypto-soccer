CREATE TABLE player_sell_orders (
    playerId NUMERIC(78, 0) NOT NULL,
    price INT NOT NULL,
    rnd NUMERIC(78,0) NOT NULL,
    validUntil NUMERIC(78,0) NOT NULL,
    typeOfTx INT NOT NULL,
    signature TEXT NOT NULL,
    PRIMARY KEY(playerId)
);

CREATE TABLE player_buy_orders (
    playerId NUMERIC(78,0) NOT NULL REFERENCES player_sell_orders(playerId),
    teamId NUMERIC(78,0) NOT NULL,
    PRIMARY KEY(playerId)
)
