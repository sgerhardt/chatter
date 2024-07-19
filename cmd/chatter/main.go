package main

import (
	"chatter/internal/client"
	"chatter/internal/config"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	run()
}

func readEnvFile() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return os.Getenv("XI_API_KEY"), os.Getenv("OUTPUT")
}

func run() {
	app, httpClient := setup()

	client.New(app, httpClient).Run()
}

func setup() (config.AppConfig, client.HTTP) {
	var app config.AppConfig
	key, dir := readEnvFile()
	if key == "" {
		log.Fatal("API Key not found")
	}
	app.APIKey = key
	app.OutputDir = dir
	app.CharacterRequestLimit = 10000

	textInput := flag.String("t", "", "Text to convert to voice")
	siteInput := flag.String("s", "", "Website to read text from")
	voiceID := flag.String("v", "", "Voice ID to use")
	flag.Parse()
	if *voiceID == "" {
		log.Fatal("Voice ID is required")
	}
	app.VoiceID = *voiceID

	app.TextInput = *textInput
	app.WebsiteURL = *siteInput
	if app.TextInput != "" && app.WebsiteURL != "" {
		log.Fatal("Only one of text or site can be provided")
	}

	httpClient := &http.Client{
		Timeout: time.Second * 310,
		Transport: &http.Transport{
			DialContext:           (&net.Dialer{Timeout: time.Second * 3}).DialContext,
			TLSHandshakeTimeout:   time.Second * 3,
			ResponseHeaderTimeout: time.Second * 300, // eleven labs doesn't appear to respond with the header until the request completes
		},
	}

	return app, httpClient
}
