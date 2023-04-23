run:
	go run -tags "static" main.go

build:
	go build -tags "static" -o build/go-istage

install:
	go install -tags "static"