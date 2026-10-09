package cmd

import (
	"fmt"

	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "AI features",
	Long: `AI features via the web/app API group.

"ai generate" streams insight generation for a card (SSE). The stream is
printed to stdout as it arrives; generation consumes AI quota.`,
}

var aiCard string

var aiGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate AI insight for a card (streams)",
	Example: `  cubox-cli ai generate --card 7247925101516031380`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if aiCard == "" {
			return fmt.Errorf("--card is required")
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		web, err := webClient(cfg)
		if err != nil {
			return err
		}
		return web.WebInsightGenerateStream(aiCard, func(line string) {
			fmt.Println(line)
		})
	},
}

func init() {
	aiCmd.AddCommand(aiGenerateCmd)
	rootCmd.AddCommand(aiCmd)

	aiGenerateCmd.Flags().StringVar(&aiCard, "card", "", "card ID (required)")
}
