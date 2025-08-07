build:
	@echo "package main" > dummy.go
	@swag init --dir . --generalInfo cmd/api/api.go --output docs
	@rm -f dummy.go
	@go build -o bin/executive cmd/main.go

run:
	@./bin/executive

swagger:
	@echo "package main" > dummy.go
	@swag init --dir . --generalInfo cmd/api/api.go --output docs
	@rm -f dummy.go
