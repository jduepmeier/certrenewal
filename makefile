.PHONY: clean

bin/cert-renewal: bin *.go cmd/cert-renewal/*.go
	go build -o bin/cert-renewal cmd/cert-renewal/*.go

bin:
	mkdir -p bin

clean:
	rm -rf bin