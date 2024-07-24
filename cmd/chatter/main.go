package main

import (
	"github.com/sgerhardt/chatter/internal/setup"
	"log"
)

func main() {
	if err := setup.NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
