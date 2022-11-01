postgres:
	docker run --name postgresql -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres
postgres-start:
	docker start postgresql
migrate:
	migrate -source file://migrations \
			-database postgres://postgres:postgres@localhost/pangea?sslmode=disable up
migrate-down:
	migrate -source file://migrations \
			-database postgres://postgres:postgres@localhost/pangea?sslmode=disable down
