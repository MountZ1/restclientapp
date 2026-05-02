package main

import (
	"log"
	"restclient/internal/tui"
)

func main() {
	err := tui.Run()
	if err != nil {
		log.Fatal(err)
	}
}
