.PHONY: build run clean

build:
	go build -o bin/fpscore.exe ./cmd

run:
	go run ./cmd

clean:
	rm -f bin/fpscore.exe

