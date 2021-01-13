.PHONY: clean test test-coverage

bin/cert-renewal: bin go.mod go.sum *.go cmd/cert-renewal/*.go
	go build -o bin/cert-renewal cmd/cert-renewal/*.go

bin:
	mkdir -p bin

clean:
	rm -rf bin

test:
	go test -cover

test-coverage: bin
	go test -coverprofile=bin/coverage.out
	go tool cover -html=bin/coverage.out
	rm bin/coverage.out
