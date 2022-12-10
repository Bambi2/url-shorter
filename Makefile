run-postgres:
	DB_TYPE=postgres docker-compose up --build url-shorter

run-redis:
	DB_TYPE=redis docker-compose up --build url-shorter 

test:
	go test -v ./...