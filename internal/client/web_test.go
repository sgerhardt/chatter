package client

import (
	"github.com/sgerhardt/chatter/internal/client/mocks"
	"github.com/sgerhardt/chatter/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestWebReader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want []string

		error     error
		charLimit int
		mockSetup func(client *mocks.HTTP)
	}{
		{
			name:      "Given a website, read the header and body",
			want:      []string{"This is the h1\nThis is paragraph text\n"},
			charLimit: 100,
			mockSetup: func(client *mocks.HTTP) {
				client.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`<html><body><h1>This is the h1</h1><p>This is paragraph text</p></body></html>`)),
				}, nil)
			},
		},
		{
			name:      "Given a website that requires batching requests",
			want:      []string{"This ", "is th", "e h1\n", "This ", "is pa", "ragra", "ph te", "xt\n"},
			charLimit: 5,
			mockSetup: func(client *mocks.HTTP) {
				client.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`<html><body><h1>This is the h1</h1><p>This is paragraph text</p></body></html>`)),
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockClient := mocks.NewHTTP(t)
			tt.mockSetup(mockClient)
			appConfig := &config.AppConfig{
				CharacterRequestLimit: tt.charLimit,
				APIKey:                "testkey",
				VoiceID:               "testvoice",
				WebsiteURL:            "https://test.com",
			}
			c := New(appConfig, mockClient)
			texts, err := c.FromWebsite("https://test.com")
			if tt.error != nil {
				assert.EqualError(t, err, tt.error.Error())
				return
			}
			assert.Equal(t, tt.want, texts)
			mockClient.AssertExpectations(t)
		})
	}
}
