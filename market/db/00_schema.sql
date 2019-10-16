CREATE TABLE auctions (
    uuid UUID NOT NULL,
    player_id TEXT NOT NULL,
    currency_id INT NOT NULL,
    price INT NOT NULL,
    rnd INT NOT NULL,
    valid_until TEXT NOT NULL,
    signature TEXT NOT NULL,
    PRIMARY KEY(uuid)
);

CREATE TABLE bids (
    auction UUID NOT NULL REFERENCES auctions(uuid),
    extra_price NUMERIC(15,2) NOT NULL,
    rnd INT NOT NULL,
    team_id TEXT NOT NULL,
    signature TEXT NOT NULL,
    PRIMARY KEY(auction, extra_price)
)
