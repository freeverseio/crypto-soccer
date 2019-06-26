FROM postgres:alpine

# WORKDIR /sql

# COPY ./sql/. /sql

ENV POSTGRES_USER freeverse
ENV POSTGRES_PASSWORD freeverse
ENV POSTGRES_DB cryptosoccer

COPY ./sql/. /docker-entrypoint-initdb.d/

