/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"time"

	fp "path/filepath"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/yendelevium/cyril/config"
	"github.com/yendelevium/cyril/internal/tui"
	bolt "go.etcd.io/bbolt"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "read a saved note, either by its name or its alias",
	Long:  "Displays the given note. You can access it via the filename or the aliasname. If there are multiple files with the same name across topics/multiple prefix-matches for that file, you can choose which one to read via a selector.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]
		aliasNames := []tui.FileData{}

		// Get all keys that have the same alias name but diff topics
		err := MatchAliasPrefixes(filename, &aliasNames)
		if err != nil {
			return err
		}

		// TODO: Walk the store to see any matching files using os.Walkdir (I think)
		// TODO: IF file doesn't exist offer to create it
		// TODO: IF multiple files exist offer to choose which one to open
		// If aliasName length is 0 or 1, handle normally
		if len(aliasNames) == 0 {
			// No alias names
			// Walk the directory? or NO?
			fmt.Printf("Couldn't find the note %s!\n", filename)
			return nil
		}

		if len(aliasNames) == 1 {
			file := aliasNames[0]
			s := tui.FileDisplay(file)
			lipgloss.Println(s)
			return nil
		}

		// Start the bubbletea program to display the options
		model := tui.ReadModel{
			Files: aliasNames,
		}

		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Couldn't run bubbletea: %v\n", err)
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(readCmd)
}

func MatchAliasPrefixes(filename string, aliasNames *[]tui.FileData) error {
	// Making this a function as I want to use defer statement so I don't forget to close the DB
	// Also using it in ReadOnly mode
	// dbPath := fmt.Sprintf("%s/cyril.db", Conf.DBPath)
	dbPath := fp.Join(config.Conf.DBPath, "cyril.db")
	db, err := bolt.Open(dbPath, 0644, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	if err != nil {
		return fmt.Errorf("Couldn't open DB: %v \n", err)
	}
	defer db.Close()

	// Get all keys with the filename as an alias
	// Stolen from official docs
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("cyril"))
		if bucket == nil {
			return fmt.Errorf("NO notes recorded yet...\n")
		}

		// Matching Prefixes (multiple topics might have same filename)
		c := bucket.Cursor()
		prefix := []byte(filename)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			*aliasNames = append(*aliasNames, tui.FileData{
				Filename: string(k),
				Filepath: string(v),
			})
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("KV error: %v\n", err)
	}
	return nil
}
