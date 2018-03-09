default: build

build:
	go fmt
	go vet -v
	go build -v

test: build
	go test -v ./...

coverage-test:
	go test -coverprofile=coverage.out -v ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out
	rm coverage.out
