package cmd

import (
	"fmt"
	"strings"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var markCmd = &cobra.Command{
	Use:   "mark",
	Short: "Manage marks (highlights)",
	Long: `Manage marks — the highlights you made while reading.

"mark list" browses highlights; "mark delete" and "mark color" manage them.
(Reading-side annotations are also exposed by "annotation list".)`,
}

var (
	markPage     int
	markKeyword  string
	markDelIDs   []string
	markColorIDs []string
	markColor    int
)

var markListCmd = &cobra.Command{
	Use:   "list",
	Short: "List marks (highlights)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		web, err := webClient(cfg)
		if err != nil {
			return err
		}
		raw, err := web.WebMarksList(markPage, markKeyword)
		if err != nil {
			return err
		}
		printJSON(jsonRaw(raw))
		return nil
	},
}

var markDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete marks (highlights)",
	Example: `  cubox-cli mark delete --id 7444025677600260245,7443973659296793971`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(markDelIDs) == 0 {
			return fmt.Errorf("--id is required")
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		web, err := webClient(cfg)
		if err != nil {
			return err
		}
		raw, err := web.WebMarksDelete(markDelIDs)
		if err != nil {
			return err
		}
		printJSON(map[string]interface{}{
			"count":   len(markDelIDs),
			"message": fmt.Sprintf("%d mark(s) deleted", len(markDelIDs)),
			"data":    jsonRaw(raw),
		})
		return nil
	},
}

var markColorCmd = &cobra.Command{
	Use:   "color",
	Short: "Update mark colors",
	Long: `Update the color of one or more marks.

colorType: 1=Yellow 2=Green 3=Blue 4=Pink 5=Purple`,
	Example: `  cubox-cli mark color --id 7444025677600260245 --color 3`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(markColorIDs) == 0 {
			return fmt.Errorf("--id is required")
		}
		if markColor < 1 || markColor > 5 {
			return fmt.Errorf("--color must be 1..5")
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		web, err := webClient(cfg)
		if err != nil {
			return err
		}
		raw, err := web.WebMarksColor(markColorIDs, markColor)
		if err != nil {
			return err
		}
		printJSON(map[string]interface{}{
			"count":   len(markColorIDs),
			"color":   markColor,
			"message": fmt.Sprintf("%d mark(s) recolored", len(markColorIDs)),
			"data":    jsonRaw(raw),
		})
		return nil
	},
}

func init() {
	markCmd.AddCommand(markListCmd)
	markCmd.AddCommand(markDeleteCmd)
	markCmd.AddCommand(markColorCmd)
	rootCmd.AddCommand(markCmd)

	markListCmd.Flags().IntVar(&markPage, "page", 1, "page number (1-based)")
	markListCmd.Flags().StringVar(&markKeyword, "keyword", "", "filter by keyword")
	markDeleteCmd.Flags().StringSliceVar(&markDelIDs, "id", nil, "mark IDs to delete (comma-separated, required)")
	markColorCmd.Flags().StringSliceVar(&markColorIDs, "id", nil, "mark IDs to recolor (comma-separated, required)")
	markColorCmd.Flags().IntVar(&markColor, "color", 0, "colorType 1..5 (required)")
}

var _ = strings.TrimSpace
var _ = client.APIResponse{}
