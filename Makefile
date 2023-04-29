.PHONY: run
run:
	go run -tags "static" main.go

.PHONY: build
build:
	go build -tags "static" -o build/go-istage

.PHONY: install
install:
	go install -tags "static"

.PHONY: test
test:
	go test -tags "static" ./...