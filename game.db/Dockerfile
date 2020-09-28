FROM postgres:12.4-alpine

ENV POSTGRES_USER freeverse
ENV POSTGRES_PASSWORD freeverse
ENV POSTGRES_DB game

COPY ./*.sql /docker-entrypoint-initdb.d/

