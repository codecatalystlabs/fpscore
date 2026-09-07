.PHONY: build run clean migrate

build:
	go build -o bin/fpscore.exe ./cmd

run:
	go run ./cmd

# Windows: migrate.bat   |  Unix/Git Bash: ./migrate.sh
migrate:
	migrate.bat

clean:
	rm -f bin/fpscore.exe

