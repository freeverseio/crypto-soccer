
CREATE TABLE unpayment (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    owner TEXT NOT NULL,
    auction_id TEXT NOT NULL,
    time_of_unpayment timestamp without time zone NOT NULL,
    notified BOOLEAN DEFAULT false,
    PRIMARY KEY(ID),
    UNIQUE(auction_id, owner)
);

CREATE INDEX idx_unpayment_owner ON unpayment (owner);
CREATE INDEX idx_unpayment_time_of_unpayment ON unpayment (time_of_unpayment);
