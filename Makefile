.env:
	cp .env.example .env


sg-tar: $(shell find . -type f -name '*.go')
	go build -o sg-tar main.go


.PHONY: clean
clean:
	rm -f sg-tar

.PHONY: lint test
lint:
	golangci-lint run

test:
	go test -v ./...
