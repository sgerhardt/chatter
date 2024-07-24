package client_test

import (
	"bytes"
	"github.com/sgerhardt/chatter/internal/client"
	"github.com/sgerhardt/chatter/internal/client/mocks"
	"github.com/sgerhardt/chatter/internal/config"
	"os"

	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestClient_ProcessText(t *testing.T) {
	t.Parallel()

	type fields struct {
		apiKey         string
		outputFilePath string
	}
	type args struct {
		text    string
		voiceID string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		error     error
		mockSetup func(client *mocks.HTTP)
	}{
		{
			name: "errors if voice id is not populated",
			args: args{
				text:    "test",
				voiceID: "",
			},
			error:     errors.New("voice ID is required"),
			mockSetup: func(_ *mocks.HTTP) {},
		},
		{
			name: "Sends text to eleven labs and writes the response to a file",
			fields: fields{
				apiKey:         "123",
				outputFilePath: t.TempDir(),
			},
			args: args{
				text:    "testing",
				voiceID: "stephen_hawking",
			},

			mockSetup: func(client *mocks.HTTP) {
				mockResp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("bytes representing the mp3 file..."))),
				}

				mockResp.StatusCode = http.StatusOK
				mockResp.Body = io.NopCloser(bytes.NewReader([]byte("bytes representing the mp3 file...")))
				// Set up the expectation
				client.On("Do", mock.AnythingOfType("*http.Request")).Return(mockResp, nil).Run(func(args mock.Arguments) {
					req := args.Get(0).(*http.Request)
					// Verify the body of the request is the expected json
					body := new(bytes.Buffer)
					_, err := body.ReadFrom(req.Body)
					require.NoError(t, err)
					assert.Equal(t, body.String(), `{"text":"testing","model_id":"eleven_monolingual_v1","voice_settings":{"stability":0,"similarity_boost":0}}`)
				})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockClient := mocks.NewHTTP(t)
			tt.mockSetup(mockClient)
			cfg := &config.AppConfig{
				CharacterRequestLimit: 100,
				OutputDir:             tt.fields.outputFilePath,
				APIKey:                tt.fields.apiKey,
				VoiceID:               tt.args.voiceID,
				TextInput:             tt.args.text,
			}
			c := client.New(cfg, mockClient)

			err := c.ProcessText()
			if tt.error != nil {
				assert.EqualError(t, err, tt.error.Error())
				return
			}

			require.NoError(t, err)
			assert.DirExists(t, tt.fields.outputFilePath)
			files, err := os.ReadDir(tt.fields.outputFilePath)
			assert.NoError(t, err)
			assert.NotEmpty(t, files, "no files found in the temporary directory")
			require.Len(t, files, 1)
			// verify the contents of the file
			file, err := os.ReadFile(tt.fields.outputFilePath + string(os.PathSeparator) + files[0].Name())
			require.NoError(t, err)
			assert.Equal(t, "bytes representing the mp3 file...", string(file))

			mockClient.AssertExpectations(t)
		})
	}
}

func TestClient_ProcessSite(t *testing.T) {
	t.Parallel()

	type fields struct {
		apiKey         string
		outputFilePath string
	}
	type args struct {
		url     string
		voiceID string
	}

	tests := []struct {
		name      string
		fields    fields
		args      args
		error     error
		mockSetup func(client *mocks.HTTP)
	}{
		{
			name: "errors if voice id is not populated",
			args: args{
				url:     "https://example.com",
				voiceID: "",
			},
			error: errors.New("voice ID is required"),
			mockSetup: func(client *mocks.HTTP) {
				mockURLResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("<html><body><h1>testing</h1></body></html>"))),
				}

				mockURLResponse.StatusCode = http.StatusOK
				mockURLResponse.Body = io.NopCloser(bytes.NewReader([]byte("bytes representing the mp3 file...")))
				mockElevenLabsResp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("bytes representing the mp3 file..."))),
				}
				mockElevenLabsResp.StatusCode = http.StatusOK
				mockElevenLabsResp.Body = io.NopCloser(bytes.NewReader([]byte("bytes representing the mp3 file...")))
				client.On("Do", mock.AnythingOfType("*http.Request")).Return(mockURLResponse, nil).Once()
			},
		},
		{
			name: "Reads text from a website, sends it to eleven labs, and writes the response to a file",
			fields: fields{
				apiKey:         "123",
				outputFilePath: t.TempDir(),
			},
			args: args{
				url:     "https://example.com",
				voiceID: "stephen_hawking",
			},

			mockSetup: func(client *mocks.HTTP) {
				// First fetch the website
				mockURLResponse := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("<html><body><h1>testing</h1></body></html>"))),
				}
				client.On("Do", mock.MatchedBy(func(req *http.Request) bool {
					return req.URL.String() == "https://example.com"
				})).Return(mockURLResponse, nil).Once()

				// Make the request to eleven labs
				mockElevenLabsResp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("bytes representing the mp3 file..."))),
				}
				mockElevenLabsResp.StatusCode = http.StatusOK
				mockElevenLabsResp.Body = io.NopCloser(bytes.NewReader([]byte("bytes representing the mp3 file...")))
				client.On("Do", mock.AnythingOfType("*http.Request")).Return(mockElevenLabsResp, nil).Run(func(args mock.Arguments) {
					req := args.Get(0).(*http.Request)
					// Verify the body of the request is the expected json
					body := new(bytes.Buffer)
					_, err := body.ReadFrom(req.Body)
					require.NoError(t, err)
					assert.Equal(t, body.String(), `{"text":"testing\n","model_id":"eleven_monolingual_v1","voice_settings":{"stability":0,"similarity_boost":0}}`)
				}).Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockClient := mocks.NewHTTP(t)
			tt.mockSetup(mockClient)
			cfg := &config.AppConfig{
				CharacterRequestLimit: 100,
				OutputDir:             tt.fields.outputFilePath,
				APIKey:                tt.fields.apiKey,
				VoiceID:               tt.args.voiceID,
				WebsiteURL:            tt.args.url,
			}
			c := client.New(cfg, mockClient)

			err := c.ProcessSite()
			if tt.error != nil {
				assert.EqualError(t, err, tt.error.Error())
				return
			}

			require.NoError(t, err)
			assert.DirExists(t, tt.fields.outputFilePath)
			files, err := os.ReadDir(tt.fields.outputFilePath)
			assert.NoError(t, err)
			assert.NotEmpty(t, files, "no files found in the temporary directory")
			require.Len(t, files, 1)
			// verify the contents of the file
			file, err := os.ReadFile(tt.fields.outputFilePath + string(os.PathSeparator) + files[0].Name())
			require.NoError(t, err)
			assert.Equal(t, "bytes representing the mp3 file...", string(file))

			mockClient.AssertExpectations(t)
		})
	}
}
