package main

import (
	"chatter/client"
	"flag"
	"log"
)

func main() {
	run()
}

func run() {
	textInput := flag.String("t", "hello world!", "Text to convert to voice")
	voiceID := flag.String("v", "", "Voice ID to use")
	flag.Parse()
	if *voiceID == "" {
		panic("Voice ID is required")
	}

	c := client.New()
	bytes, err := c.GenerateVoiceForText(*textInput, *voiceID)
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}
