run-postgres:
	DB_TYPE=postgres docker-compose up --build url-shorter

run-redis:
	DB_TYPE=redis docker-compose up --build url-shorter 

postgres-run:
	docker run --name url-shorter -e POSTGRES_PASSWORD=qwerty -p 5436:5432 -d --rm postgres

redis-run:
	docker run --name url-shorter-redis -p 6379:6379 -d --rm redis

test:
	go test -v ./...