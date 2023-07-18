run:
	@go run main.go

test:
	@go test -v ./...

seed:
	@go run scripts/seed.go