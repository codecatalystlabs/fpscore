.PHONY: build run clean migrate

# Linux/Ubuntu binary name; on Windows use: go build -o bin/fpscore.exe ./cmd
build:
	mkdir -p bin
	go build -o bin/fpscore ./cmd

run:
	go run ./cmd

# Apply SQL migrations (Ubuntu/Linux). On Windows run migrate.bat instead.
migrate:
	chmod +x migrate.sh
	./migrate.sh

clean:
	rm -f bin/fpscore bin/fpscore.exe
