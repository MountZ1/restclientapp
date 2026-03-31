.PHONY: build run dev clean

build:
	go build -o ./tmp/tui.exe ./cmd/tui

run:
	./tmp/tui.exe

clean:
	rm -rf ./tmp/tui.exe

dev:	build run
