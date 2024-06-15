package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ElvenRequest struct {
	Text                            string                            `json:"text"`
	ModelId                         string                            `json:"model_id"`
	VoiceSettings                   VoiceSettings                     `json:"voice_settings"`
	PronunciationDictionaryLocators []PronunciationDictionaryLocators `json:"pronunciation_dictionary_locators,omitempty"`
	Seed                            int                               `json:"seed,omitempty"`
	PreviousText                    string                            `json:"previous_text,omitempty"`
	NextText                        string                            `json:"next_text,omitempty"`
	PreviousRequestIds              []string                          `json:"previous_request_ids,omitempty"`
	NextRequestIds                  []string                          `json:"next_request_ids,omitempty"`
}

type VoiceSettings struct {
	Stability       float64 `json:"stability"`
	SimilarityBoost float64 `json:"similarity_boost"`
	Style           int     `json:"style,omitempty"`
	UseSpeakerBoost bool    `json:"use_speaker_boost,omitempty"`
}

type PronunciationDictionaryLocators struct {
	PronunciationDictionaryId string `json:"pronunciation_dictionary_id,omitempty"`
	VersionId                 string `json:"version_id,omitempty"`
}

type Client struct{}

var key string
var output string

func init() {
	key, output = readEnvFile()
	if key == "" {
		log.Fatal("API Key not found")
	}
	if output == "" {
		currentTime := time.Now()
		formattedTime := currentTime.Format("20060102_150405")
		output = "output_" + fmt.Sprintf("%s", formattedTime) + ".mp3"
	}
}

func New() *Client {
	return &Client{}
}

func (c *Client) GenerateVoiceFromText(text string, voiceID string) {
	payload, err := buildPayload(text)
	if err != nil {
		panic(err)
	}

	req := buildRequest(err, voiceID, payload)
	body := doRequest(err, req)
	writeFile(err, body)
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
		ModelId: "eleven_monolingual_v1",
		VoiceSettings: VoiceSettings{
			Stability:       0,
			SimilarityBoost: 0,
			Style:           0,
			UseSpeakerBoost: false,
		},
		PronunciationDictionaryLocators: nil,
		Seed:                            0,
		PreviousText:                    "",
		NextText:                        "",
		PreviousRequestIds:              nil,
		NextRequestIds:                  nil,
	}

	payload, err := json.Marshal(elvenReq)
	if err != nil {
		panic(err)
	}
	return payload, err
}

func writeFile(err error, body []byte) {
	err = os.WriteFile(output, body, 0644)
	fmt.Println("File written to", output)
}

func doRequest(err error, req *http.Request) []byte {
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		cErr := Body.Close()
		if cErr != nil {
			panic(cErr)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func buildRequest(err error, voiceID string, payload []byte) *http.Request {
	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "audio/mpeg")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("xi-api-key", key)
	return req
}
