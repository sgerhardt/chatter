package main

import (
	"github.com/sgerhardt/chatter/internal/client"
	"github.com/sgerhardt/chatter/internal/setup"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	run()
}

func run() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "an eleven labs client for text to voice",
	Run: func(_ *cobra.Command, _ []string) {
		cfg, c, err := setup.New(".env", voiceID, textInput, siteInput)
		if err != nil {
			log.Fatal(err)
		}
		client.New(cfg, c).Run()
	},
}
var (
	voiceID   string
	textInput string
	siteInput string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&textInput, "text", "t", "", "Text to convert to voice")
	rootCmd.PersistentFlags().StringVarP(&siteInput, "site", "s", "", "Website to read text from")
	rootCmd.PersistentFlags().StringVarP(&voiceID, "voice", "v", "", "Voice ID to use")
}
