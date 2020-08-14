echo y | docker-compose -f ./docker-compose_market.yml rm
docker-compose -f ./docker-compose_market.yml build --no-cache
docker-compose -f ./docker-compose_market.yml up --force-recreate
