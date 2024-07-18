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
	textInput := flag.String("t", "", "Text to convert to voice")
	siteInput := flag.String("s", "", "Website to read text from")
	voiceID := flag.String("v", "", "Voice ID to use")
	flag.Parse()
	if *voiceID == "" {
		log.Fatal("Voice ID is required")
	}

	c := client.New()

	if textInput == nil && siteInput == nil {
		log.Fatal("Either text or site must be provided")

	}

	if textInput != nil && *textInput != "" && siteInput != nil && *siteInput != "" {
		log.Fatal("Only one of text or site can be provided")
	}

	if textInput != nil && *textInput != "" {
		fromText(c, textInput, voiceID)
	}

	if siteInput != nil && *siteInput != "" {
		fromSite(c, siteInput, voiceID)
	}

}

func fromText(c *client.Client, textInput *string, voiceID *string) {
	bytes, err := c.FromText(*textInput, *voiceID)
	if err != nil {
		log.Fatal(err)
	}
	_, err = c.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}
}

func fromSite(c *client.Client, siteInput *string, voiceID *string) {
	texts, err := c.FromWebsite(*siteInput)
	if err != nil {
		log.Fatal(err)
	}
	for _, text := range texts {
		bytes, tErr := c.FromText(text, *voiceID)
		if tErr != nil {
			log.Fatal(tErr)
		}
		_, err = c.Write(bytes)
		if err != nil {
			log.Fatal(err)
		}
	}

}
