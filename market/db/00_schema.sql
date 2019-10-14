CREATE TABLE player_sell_orders (
    player_id TEXT NOT NULL,
    currency_id INT NOT NULL,
    price INT NOT NULL,
    rnd INT NOT NULL,
    valid_until TEXT NOT NULL,
    signature TEXT NOT NULL,
    PRIMARY KEY(player_id)
);

CREATE TABLE player_buy_orders (
    player_id TEXT NOT NULL,
    extra_price NUMERIC(15,2) NOT NULL,
    rnd INT NOT NULL,
    team_id TEXT NOT NULL,
    is_2_start_auction BOOLEAN NOT NULL,
    signature TEXT NOT NULL,
    PRIMARY KEY(player_id)
)
