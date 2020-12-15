
CREATE TABLE unpayment (
    owner TEXT NOT NULL,
    num_of_unpayments INT DEFAULT 0,
    last_time_of_unpayment timestamp without time zone NOT NULL,
    PRIMARY KEY(owner)
);
