package config

// AppConfig holds the application config - it should not import any other packages

type AppConfig struct {
	CharacterRequestLimit int
	TextInput             string
	OutputDir             string
	APIKey                string
	VoiceID               string
	WebsiteURL            string
}
