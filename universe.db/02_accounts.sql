CREATE TABLE accounts(
    address TEXT NOT NULL,
    PRIMARY KEY (address)
);
comment on table accounts is E'@omit create,delete';
comment on column accounts.address is E'@omit update';

CREATE TABLE accounts_histories(
    owner TEXT NOT NULL,
    PRIMARY KEY (owner)
);
comment on table accounts_histories is E'@omit create,update,delete';