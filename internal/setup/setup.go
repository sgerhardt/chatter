package setup

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sgerhardt/chatter/internal/client"
	"github.com/sgerhardt/chatter/internal/config"
	"github.com/spf13/cobra"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func NewRootCmd() *cobra.Command {

	var voiceID string
	var textInput string
	var siteInput string

	cmd := &cobra.Command{
		Use:   "chatter -v <voiceID> {-t <text> | -s <url>}",
		Short: "An Eleven Labs client for text to voice",
		Long: `Chatter is a command-line client for Eleven Labs text-to-voice service.

Usage:
  chatter -v <voiceID> -t <text>   (Provide text to convert to voice)
  chatter -v <voiceID> -s <url>    (Provide a URL to read text from)

Either --text or --site is required, but not both.`,
		PreRunE: func(_ *cobra.Command, _ []string) error {
			if voiceID == "" {
				return errors.New("voice is required")
			}
			if textInput == "" && siteInput == "" {
				return errors.New("text or site is required")
			}
			if textInput != "" && siteInput != "" {
				return errors.New("only one of text or site can be provided")
			}
			return nil
		},
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, c, err := New(".env", voiceID, textInput, siteInput)
			if err != nil {
				return err
			}
			if textInput != "" {
				return client.New(cfg, c).ProcessText()
			} else if siteInput != "" {
				return client.New(cfg, c).ProcessSite()
			}
			return errors.New("text or site is required")
		},
	}

	cmd.Flags().StringVarP(&textInput, "text", "t", "", "Text to convert to voice")
	cmd.Flags().StringVarP(&siteInput, "site", "s", "", "Website to read text from")
	cmd.Flags().StringVarP(&voiceID, "voice", "v", "", "Voice ID to use")
	if err := cmd.MarkFlagRequired("voice"); err != nil {
		log.Fatal(err)
	}

	return cmd
}

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
