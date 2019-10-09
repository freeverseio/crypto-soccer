FROM postgres:11.5-alpine

ENV POSTGRES_USER freeverse
ENV POSTGRES_PASSWORD freeverse
ENV POSTGRES_DB market

COPY ./*.sql /docker-entrypoint-initdb.d/

