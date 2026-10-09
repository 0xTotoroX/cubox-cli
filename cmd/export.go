package cmd

import (
	"fmt"

	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export your library",
	Long: `Export bookmarks via the web/app API group.

"export status" shows today's export count and account info;
"export bookmarks" requests a full-library export (the server responds with
the export task payload; large libraries are processed asynchronously).`,
}

var exportStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show export-related counters",
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
		if raw, err := web.WebGetRaw("/c/api/bookmark/export/today/count"); err == nil {
			out["todayCount"] = jsonRaw(raw)
		} else {
			out["todayCount_error"] = err.Error()
		}
		if raw, err := web.WebGetRaw("/c/api/import/ready"); err == nil {
			out["importReady"] = jsonRaw(raw)
		}
		printJSON(out)
		return nil
	},
}

var exportBookmarksCmd = &cobra.Command{
	Use:   "bookmarks",
	Short: "Request a full-library bookmark export",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		web, err := webClient(cfg)
		if err != nil {
			return err
		}
		raw, err := web.WebGetRaw("/c/api/v3/bookmark/export")
		if err != nil {
			return err
		}
		printJSON(map[string]interface{}{
			"message": "export requested — see data for the task/URL payload",
			"data":    jsonRaw(raw),
		})
		return nil
	},
}

func init() {
	exportCmd.AddCommand(exportStatusCmd)
	exportCmd.AddCommand(exportBookmarksCmd)
	rootCmd.AddCommand(exportCmd)
}

var _ = fmt.Sprintf
