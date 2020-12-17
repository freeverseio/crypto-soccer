
CREATE TABLE unpayment (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    owner TEXT NOT NULL,
    last_time_of_unpayment timestamp without time zone NOT NULL,
    notified BOOLEAN DEFAULT false,
    PRIMARY KEY(ID)
);
