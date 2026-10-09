package cmd

import (
	"encoding/json"
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

var settingsJSON string

var accountSettingsUpdateCmd = &cobra.Command{
	Use:   "settings-update",
	Short: "Write reading settings (experimental)",
	Long: `Write reading settings back via POST /c/api/settings/read/update.

Experimental: the payload shape is not fully documented — pass the FULL
settings object as JSON (take the output of "account settings" as the base
and modify the fields you want). Values sent are applied as-is.`,
	Example: `  cubox-cli account settings-update --json '{"markAsReadConfig":1}'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if settingsJSON == "" {
			return fmt.Errorf("--json is required (full settings object)")
		}
		var body map[string]interface{}
		if err := json.Unmarshal([]byte(settingsJSON), &body); err != nil {
			return fmt.Errorf("parsing --json: %w", err)
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		web, err := webClient(cfg)
		if err != nil {
			return err
		}
		raw, err := web.WebSettingsUpdate(body)
		if err != nil {
			return err
		}
		printJSON(map[string]interface{}{"message": "settings updated", "data": jsonRaw(raw)})
		return nil
	},
}

func init() {
	accountCmd.AddCommand(accountApikeyCmd)
	accountCmd.AddCommand(accountSettingsCmd)
	accountCmd.AddCommand(accountSettingsUpdateCmd)
	accountCmd.AddCommand(accountSyncCmd)
	accountCmd.AddCommand(accountInsightCmd)
	rootCmd.AddCommand(accountCmd)

	accountInsightCmd.Flags().StringVar(&insightCardID, "card", "", "card ID (required)")
	accountSettingsUpdateCmd.Flags().StringVar(&settingsJSON, "json", "", "full settings object as JSON (required)")
}
