package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sgerhardt/chatter/internal/config"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

type elevenRequest struct {
	Text                            string                            `json:"text"`
	ModelID                         string                            `json:"model_id"`
	VoiceSettings                   voiceSettings                     `json:"voice_settings"`
	PronunciationDictionaryLocators []pronunciationDictionaryLocators `json:"pronunciation_dictionary_locators,omitempty"`
	Seed                            int                               `json:"seed,omitempty"`
	PreviousText                    string                            `json:"previous_text,omitempty"`
	NextText                        string                            `json:"next_text,omitempty"`
	PreviousRequestIDs              []string                          `json:"previous_request_ids,omitempty"`
	NextRequestIDs                  []string                          `json:"next_request_ids,omitempty"`
}

type voiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           int     `json:"style,omitempty"`
	UseSpeakerBoost bool    `json:"use_speaker_boost,omitempty"`
}

type pronunciationDictionaryLocators struct {
	PronunciationDictionaryID string `json:"pronunciation_dictionary_id,omitempty"`
	VersionID                 string `json:"version_id,omitempty"`
}

type ElevenLabs struct {
	httpClient HTTP
	Config     *config.AppConfig
}

type HTTP interface {
	Do(req *http.Request) (*http.Response, error)
}

func (c *ElevenLabs) fileWithTimestamp() string {
	currentTime := time.Now()
	formattedTime := currentTime.Format("20060102_150405")
	prefix := ""
	if c.Config.OutputDir == "" {
		prefix = "output_"
	} else if !strings.HasSuffix(c.Config.OutputDir, string(os.PathSeparator)) {
		prefix = c.Config.OutputDir + string(os.PathSeparator)
	}
	return prefix + formattedTime + ".mp3"
}

func (c *ElevenLabs) write(data []byte) (int, error) {
	err := os.WriteFile(c.fileWithTimestamp(), data, 0644)
	if err != nil {
		return 0, err
	}

	return len(data), nil
}

func New(cfg *config.AppConfig, httpClient HTTP) *ElevenLabs {
	return &ElevenLabs{
		Config:     cfg,
		httpClient: httpClient,
	}
}

func (c *ElevenLabs) ProcessText() error {
	fromText, err := c.FromText(c.Config.TextInput, c.Config.VoiceID)
	if err != nil {
		return err
	}
	_, err = c.write(fromText)
	return err
}

func (c *ElevenLabs) ProcessSite() error {
	texts, err := c.FromWebsite(c.Config.WebsiteURL)
	if err != nil {
		return err
	}
	for _, text := range texts {
		fromText, tErr := c.FromText(text, c.Config.VoiceID)
		if tErr != nil {
			return tErr
		}
		_, err = c.write(fromText)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ElevenLabs) FromText(text string, voiceID string) ([]byte, error) {
	if count := utf8.RuneCountInString(text); count > c.Config.CharacterRequestLimit {
		return nil, fmt.Errorf("text limit is %d characters, got :%d", c.Config.CharacterRequestLimit, count)
	}
	if voiceID == "" {
		return nil, fmt.Errorf("voice ID is required")
	}

	payload, err := buildPayload(text)
	if err != nil {
		return nil, fmt.Errorf("failed to build payload: %w", err)
	}

	req, err := buildRequest(c.Config.APIKey, voiceID, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func buildPayload(text string) ([]byte, error) {
	elvenReq := elevenRequest{
		Text:    text,
		ModelID: "eleven_monolingual_v1",
		VoiceSettings: voiceSettings{
			Stability:       0,
			SimilarityBoost: 0,
		},
	}
	return json.Marshal(elvenReq)
}

func (c *ElevenLabs) doRequest(req *http.Request) ([]byte, error) {
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
