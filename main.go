package main

import (
	"chatter/client"
	"flag"
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
	c.GenerateVoiceFromText(*textInput, *voiceID)
}
