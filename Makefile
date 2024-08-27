run:
	go run ./cmd/api

sqlcgen:
	sqlc generate

createMigrate:
	@echo "Creating migration: $(name)"
	@if [ -z "$(name)" ]; then \
    	    migrate create -ext=sql -dir=sql/migrations -seq init; \
    	else \
    	    migrate create -ext=sql -dir=sql/migrations -seq $(name); \
    	fi
migrateUp:
	migrate -path=sql/migrations -database "postgres://postgres:postgres@localhost/guapi_teste?sslmode=disable" -verbose up

migrateDown:
	migrate -path=sql/migrations -database "postgres://postgres:postgres@localhost/sistema_de_passagem?sslmode=disable" -verbose down
