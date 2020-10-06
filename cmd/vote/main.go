package main

import (
	"log"

	"github.com/copilot-example-voting-app/vote"
)

func main() {
	if err := vote.Run(); err != nil {
		log.Fatalf("run vote server: %v\n", err)
	}
}
