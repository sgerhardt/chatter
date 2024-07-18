package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

const RequestCharacterLimit = 10000

type ElvenRequest struct {
	Text                            string                            `json:"text"`
	ModelID                         string                            `json:"model_id"`
	VoiceSettings                   VoiceSettings                     `json:"voice_settings"`
	PronunciationDictionaryLocators []PronunciationDictionaryLocators `json:"pronunciation_dictionary_locators,omitempty"`
	Seed                            int                               `json:"seed,omitempty"`
	PreviousText                    string                            `json:"previous_text,omitempty"`
	NextText                        string                            `json:"next_text,omitempty"`
	PreviousRequestIDs              []string                          `json:"previous_request_ids,omitempty"`
	NextRequestIDs                  []string                          `json:"next_request_ids,omitempty"`
}

type VoiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           int     `json:"style,omitempty"`
	UseSpeakerBoost bool    `json:"use_speaker_boost,omitempty"`
}

type PronunciationDictionaryLocators struct {
	PronunciationDictionaryID string `json:"pronunciation_dictionary_id,omitempty"`
	VersionID                 string `json:"version_id,omitempty"`
}

type Client struct {
	httpClient     HttpClient
	apiKey         string
	outputFilePath string
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (c *Client) fileWithTimestamp() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("20060102_150405")
	prefix := ""
	if c.outputFilePath == "" {
		prefix = "output_"
	} else if !strings.HasSuffix(c.outputFilePath, string(os.PathSeparator)) {
		prefix = c.outputFilePath + string(os.PathSeparator)
	}
	return prefix + fmt.Sprintf("%s", formattedTime) + ".mp3"
}

func (c *Client) Write(data []byte) (int, error) {
	err := os.WriteFile(c.fileWithTimestamp(), data, 0644)
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

func New() *Client {
	apiKey, output := readEnvFile()
	if apiKey == "" {
		log.Fatal("API Key not found")
	}
	return &Client{
		apiKey:         apiKey,
		outputFilePath: output,
		httpClient: &http.Client{
			Timeout: time.Second * 310,
			Transport: &http.Transport{
				DialContext:           (&net.Dialer{Timeout: time.Second * 3}).DialContext,
				TLSHandshakeTimeout:   time.Second * 3,
				ResponseHeaderTimeout: time.Second * 300, // eleven labs doesn't appear to respond with the header until the request completes
			},
		},
	}
}

func (c *Client) FromText(text string, voiceID string) ([]byte, error) {
	if count := utf8.RuneCountInString(text); count > RequestCharacterLimit {
		return nil, fmt.Errorf("text limit is %d characters, got :%d", RequestCharacterLimit, count)
	}
	if voiceID == "" {
		return nil, fmt.Errorf("voice ID is required")
	}

	payload, err := buildPayload(text)
	if err != nil {
		return nil, fmt.Errorf("failed to build payload: %w", err)
	}

	req, err := buildRequest(c.apiKey, voiceID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func readEnvFile() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	return os.Getenv("XI_API_KEY"), os.Getenv("OUTPUT")
}

func buildPayload(text string) ([]byte, error) {
	elvenReq := ElvenRequest{
		Text:    text,
		ModelID: "eleven_monolingual_v1",
		VoiceSettings: VoiceSettings{
			Stability:       0,
			SimilarityBoost: 0,
		},
	}
	return json.Marshal(elvenReq)
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Printf("error closing response body: %v", closeErr)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed: %s, body:%v", res.Status, string(body))
	}
	return body, nil
}

func buildRequest(apiKey, voiceID string, payload []byte) (*http.Request, error) {
	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "audio/mpeg")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("xi-api-key", apiKey)
	return req, nil
}
