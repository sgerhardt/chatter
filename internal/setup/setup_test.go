package setup

import (
	"github.com/sgerhardt/chatter/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// nolint:paralleltest
	// This test deals with setting os-level env vars, which is not supported in parallel tests
	type expected struct {
		errString string
		cfg       *config.AppConfig
	}

	tests := []struct {
		name string

		envFile string

		expected expected
	}{
		{
			name: "errors when .env file not present",
			expected: expected{
				errString: ".env file not found",
				cfg:       &config.AppConfig{},
			},
			envFile: "",
		},
		{
			name: ".env file exists and is empty",
			expected: expected{
				errString: "API Key not found",
				cfg:       &config.AppConfig{},
			},
			envFile: ".env",
		},
		{
			name: "API key is in .env file and voice ID is populated",
			expected: expected{
				errString: "text or site is required",
				cfg:       &config.AppConfig{APIKey: "123", VoiceID: "testVoiceID"},
			},
			envFile: ".env",
		},
		{
			name: "API key is in .env file and voice ID is populated and text input is populated",
			expected: expected{
				errString: "",
				cfg:       &config.AppConfig{APIKey: "123", VoiceID: "testVoiceID", TextInput: "hello world!", CharacterRequestLimit: 10000},
			},
			envFile: ".env",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// nolint:paralleltest
			// This test deals with setting os-level env vars, which is not supported in parallel tests

			dir := t.TempDir()
			envFile := ""
			if tt.envFile != "" {
				file, err := os.CreateTemp(dir, "*"+tt.envFile)
				require.NoError(t, err)
				_, err = file.WriteString("XI_API_KEY=" + tt.expected.cfg.APIKey)
				t.Setenv("XI_API_KEY", tt.expected.cfg.APIKey)
				require.NoError(t, file.Close())
				require.NoError(t, err)
				t.Cleanup(func() {
					require.NoError(t, os.Remove(file.Name()))
				})
				envFile = file.Name()
			}

			// Run test
			cfg, client, err := New(envFile, tt.expected.cfg.VoiceID, tt.expected.cfg.TextInput, tt.expected.cfg.WebsiteURL)

			// Assert expectations
			if tt.expected.errString != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expected.errString)
				assert.Nil(t, cfg)
				assert.Nil(t, client)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, cfg)
			assert.NotNil(t, client)
			assert.Equal(t, tt.expected.cfg, cfg)
		})
	}
}
