build:
	@swag init --generalInfo cmd/api/api.go --output docs
	@go build -o bin/executive cmd/main.go

run:
	@./bin/executive

swagger:
	@swag init --generalInfo cmd/api/api.go --output docs
