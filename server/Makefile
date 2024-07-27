build:
	go build github.com/Prokopevs/mini/server/cmd/server

migrateup:
	migrate -path db/migrations -database "postgres://postgres:postgres@172.20.0.2:5432/minidb?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgres://postgres:postgres@172.20.0.2:5432/minidb?sslmode=disable" -verbose down