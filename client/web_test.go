package client

import (
	"chatter/client/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestWebReader(t *testing.T) {
	type fields struct {
		apiKey         string
		outputFilePath string
	}
	type args struct {
		text    string
		voiceID string
	}
	tests := []struct {
		fields    fields
		args      args
		name      string
		want      []string
		error     error
		mockSetup func(client *mocks.HttpClient)
	}{
		{
			name: "Given a website, read the header and body",
			want: []string{"This is the h1\nThis is paragraph text\n"},
			mockSetup: func(client *mocks.HttpClient) {
				client.On("Do", mock.Anything).Return(&http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`<html><body><h1>This is the h1</h1><p>This is paragraph text</p></body></html>`)),
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := mocks.NewHttpClient(t)
			tt.mockSetup(mockClient)
			c := &Client{
				apiKey:         tt.fields.apiKey,
				outputFilePath: tt.fields.outputFilePath,
				httpClient:     mockClient,
			}
			texts, err := c.FromWebsite("https://test.com")
			if tt.error != nil {
				assert.EqualError(t, err, tt.error.Error())
				return
			}
			assert.Equal(t, tt.want, texts)
		})
	}
}
