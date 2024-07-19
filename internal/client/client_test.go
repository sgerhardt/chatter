package client

import (
	"bytes"
	"chatter/internal/client/mocks"
	"chatter/internal/config"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestClient_GenerateVoiceFromText(t *testing.T) {
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
		mockSetup func(client *mocks.Http)
	}{
		{
			name: "errors if voice id is not populated",
			args: args{
				text:    "test",
				voiceID: "",
			},
			error:     errors.New("voice ID is required"),
			mockSetup: func(_ *mocks.Http) {},
		},
		{
			name:   "marshals a payload to json",
			fields: fields{},
			args: args{
				text:    "testing",
				voiceID: "stephen_hawking",
			},

			mockSetup: func(client *mocks.Http) {
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
			mockClient := mocks.NewHttp(t)
			tt.mockSetup(mockClient)
			cfg := config.AppConfig{
				CharacterRequestLimit: 100,
				OutputDir:             tt.fields.outputFilePath,
				APIKey:                tt.fields.apiKey,
				VoiceID:               tt.args.voiceID,
			}
			c := New(cfg, mockClient)
			_, err := c.FromText(tt.args.text, tt.args.voiceID)
			if tt.error != nil {
				assert.EqualError(t, err, tt.error.Error())
				return
			}

			require.NoError(t, err)

			mockClient.AssertExpectations(t)
		})
	}
}
