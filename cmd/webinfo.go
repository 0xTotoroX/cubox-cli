package cmd

import (
	"fmt"

	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Account info, settings, sync status, API keys",
	Long: `Read-only views into account-level state via the web/app API group:

  account apikey     show the API extension key(s)
  account settings   show reading settings
  account sync       show Notion / Flowus / Readwise sync status
  account insight    show a card's AI insight (by card id)`,
}

var accountApikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Show API extension key(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		web, err := webClient(cfg)
		if err != nil {
			return err
		}
		raw, err := web.WebGetRaw("/c/api/user/apiKey")
		if err != nil {
			return err
		}
		printJSON(jsonRaw(raw))
		return nil
	},
}

var accountSettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Show reading settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		web, err := webClient(cfg)
		if err != nil {
			return err
		}
		raw, err := web.WebGetRaw("/c/api/settings/read")
		if err != nil {
			return err
		}
		printJSON(jsonRaw(raw))
		return nil
	},
}

var accountSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Show Notion / Flowus / Readwise sync status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		web, err := webClient(cfg)
		if err != nil {
			return err
		}
		out := map[string]interface{}{}
		for name, path := range map[string]string{
			"notion":  "/c/api/notionSync/info",
			"flowus":  "/c/api/flowusSync/info",
			"readwise": "/c/api/readwiseSync/info",
		} {
			if raw, err := web.WebGetRaw(path); err == nil {
				out[name] = jsonRaw(raw)
			} else {
				out[name] = map[string]string{"error": err.Error()}
			}
		}
		printJSON(out)
		return nil
	},
}

var insightCardID string

var accountInsightCmd = &cobra.Command{
	Use:   "insight",
	Short: "Show a card's AI insight",
	Example: `  cubox-cli account insight --card 7247925101516031380`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if insightCardID == "" {
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
		raw, err := web.WebGetRaw("/c/api/card/insight/" + insightCardID)
		if err != nil {
			return err
		}
		printJSON(jsonRaw(raw))
		return nil
	},
}

func init() {
	accountCmd.AddCommand(accountApikeyCmd)
	accountCmd.AddCommand(accountSettingsCmd)
	accountCmd.AddCommand(accountSyncCmd)
	accountCmd.AddCommand(accountInsightCmd)
	rootCmd.AddCommand(accountCmd)

	accountInsightCmd.Flags().StringVar(&insightCardID, "card", "", "card ID (required)")
}
