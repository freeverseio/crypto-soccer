.PHONY: test docker-up docker-down 

test:
	@PGPASSWORD=freeverse psql cryptosoccer -U freeverse -h localhost -p 5432 < ./sql/schema.sql
	@go test -v ./...

docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down --volumes
