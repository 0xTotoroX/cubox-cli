package cmd

import (
	"fmt"

	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var listsCmd = &cobra.Command{
	Use:   "lists",
	Short: "Manage reading lists",
	Long: `Manage reading lists — curated, shareable collections of cards.

Note: creating/deleting reading lists themselves is available in the app;
this command group covers listing them and adding/removing cards.`,
}

var listsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List your reading lists",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		web, err := webClient(cfg)
		if err != nil {
			return err
		}
		raw, err := web.WebReadingLists()
		if err != nil {
			return err
		}
		printJSON(jsonRaw(raw))
		return nil
	},
}

var (
	addListID  string
	addCardID  string
	addHl      bool
	addNote    bool
	rmListID   string
	rmCardID   string
)

var listsAddItemCmd = &cobra.Command{
	Use:   "add-item",
	Short: "Add a card to a reading list",
	Example: `  cubox-cli lists add-item --list 7123... --card 7456... --highlight --note`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if addListID == "" || addCardID == "" {
			return fmt.Errorf("--list and --card are required")
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		web, err := webClient(cfg)
		if err != nil {
			return err
		}
		raw, err := web.WebReadingListAddItem(addListID, addCardID, addHl, addNote)
		if err != nil {
			return err
		}
		printJSON(map[string]interface{}{"message": "card added to reading list", "data": jsonRaw(raw)})
		return nil
	},
}

var listsRemoveItemCmd = &cobra.Command{
	Use:   "remove-item",
	Short: "Remove a card from a reading list",
	Example: `  cubox-cli lists remove-item --list 7123... --card 7456...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if rmListID == "" || rmCardID == "" {
			return fmt.Errorf("--list and --card are required")
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		web, err := webClient(cfg)
		if err != nil {
			return err
		}
		raw, err := web.WebReadingListRemoveItem(rmListID, rmCardID)
		if err != nil {
			return err
		}
		printJSON(map[string]interface{}{"message": "card removed from reading list", "data": jsonRaw(raw)})
		return nil
	},
}

func init() {
	listsCmd.AddCommand(listsListCmd)
	listsCmd.AddCommand(listsAddItemCmd)
	listsCmd.AddCommand(listsRemoveItemCmd)
	rootCmd.AddCommand(listsCmd)

	listsAddItemCmd.Flags().StringVar(&addListID, "list", "", "reading list ID (required)")
	listsAddItemCmd.Flags().StringVar(&addCardID, "card", "", "card ID (required)")
	listsAddItemCmd.Flags().BoolVar(&addHl, "highlight", false, "include highlights")
	listsAddItemCmd.Flags().BoolVar(&addNote, "note", false, "include notes")
	listsRemoveItemCmd.Flags().StringVar(&rmListID, "list", "", "reading list ID (required)")
	listsRemoveItemCmd.Flags().StringVar(&rmCardID, "card", "", "card ID (required)")
}
