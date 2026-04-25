.PHONY: build run dev clean

build-tui:
	go build -o ./tmp/tui.exe ./cmd/tui

build-gui:
	go build -o ./tmp/gui.exe ./cmd/gui

run-tui:
	./tmp/tui.exe

run-gui:
	./tmp/gui.exe

clean-tui:
	rm -rf ./tmp/tui.exe

clean-gui:
	rm -rf ./tmp/gui.exe

dev-gui:	build-gui run-gui
dev-tui: build-tui run-tui
