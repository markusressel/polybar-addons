OUTPUT_DIR=bin/

build: clean
	go build -o bin/battery ./cmd/battery/main.go
	go build -o bin/disk ./cmd/disk/main.go
	go build -o bin/network ./cmd/network/main.go
	go build -o bin/zfs ./cmd/zfs/main.go

deploy:	build
	cp ./bin/* /home/markus/.config/polybar/scripts/

clean:
	go clean
	rm -rf bin