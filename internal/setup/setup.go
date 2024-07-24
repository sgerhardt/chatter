package setup

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sgerhardt/chatter/internal/client"
	"github.com/sgerhardt/chatter/internal/config"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func readEnvFile(filename string) (string, string, error) {
	err := godotenv.Load(filename)
	if err != nil {
		return "", "", fmt.Errorf("error loading .env file: %v", err)
	}
	return os.Getenv("XI_API_KEY"), os.Getenv("OUTPUT"), nil
}

func New(filename string, voiceID string, textInput string, siteInput string) (*config.AppConfig, client.HTTP, error) {
	if filename == "" || !strings.HasSuffix(filename, ".env") {
		return nil, nil, errors.New(".env file not found")
	}
	key, dir, err := readEnvFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading env file: %w", err)
	}
	if key == "" {
		return nil, nil, fmt.Errorf("API Key not found")
	}

	app := &config.AppConfig{}
	app.APIKey = key
	app.OutputDir = dir
	app.CharacterRequestLimit = 10000

	if voiceID == "" {
		return nil, nil, errors.New("voice ID is required")
	}
	app.VoiceID = voiceID

	if textInput == "" && siteInput == "" {
		return nil, nil, errors.New("text or site is required")
	}
	if textInput != "" && siteInput != "" {
		return nil, nil, errors.New("only one of text or site can be provided")
	}
	app.TextInput = textInput
	app.WebsiteURL = siteInput

	httpClient := &http.Client{
		Timeout: time.Second * 310,
		Transport: &http.Transport{
			DialContext:           (&net.Dialer{Timeout: time.Second * 3}).DialContext,
			TLSHandshakeTimeout:   time.Second * 3,
			ResponseHeaderTimeout: time.Second * 300, // eleven labs doesn't appear to respond with the header until the request completes
		},
	}

	return app, httpClient, nil
}
