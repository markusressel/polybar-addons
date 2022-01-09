OUTPUT_DIR=bin/

build:
	go build -o bin/battery ./cmd/battery/main.go

clean:
	go clean
	rm bin