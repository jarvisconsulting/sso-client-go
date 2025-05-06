run:
	go run cmd/server/*.go

build:
	go build -o sso-go-client cmd/server/*.go

run-build:
	./sso-go-client
