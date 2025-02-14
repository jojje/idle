VERSION := $(shell git describe --dirty --tags)
BUILD = go build -ldflags "-X main.version=$(VERSION)"

build:
	$(BUILD)

test:
	go test ./...

build-release:
	mkdir -p dist
	GOOS=linux $(BUILD) -o dist/idle_linux
	GOOS=darwin $(BUILD) -o dist/idle_macos
	GOOS=windows $(BUILD) -o dist/idle_windows.exe

clean:
	rm -rf dist
