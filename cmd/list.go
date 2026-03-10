/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	fp "path/filepath"
	"sort"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
	"github.com/yendelevium/cyril/config"
	"github.com/yendelevium/cyril/internal/tui"
	bolt "go.etcd.io/bbolt"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all files (total and/or under a specific topic)",
	Long:  "Lists available files system wide or under a specific topic.",
	RunE: func(cmd *cobra.Command, args []string) error {
		topic, err := cmd.Flags().GetString("topic")
		if err != nil {
			return err
		}

		dbPath := fp.Join(config.Conf.DBPath, "cyril.db")
		db, err := bolt.Open(dbPath, 0644, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
		if err != nil {
			return fmt.Errorf("Couldn't open DB: %v \n", err)
		}
		defer db.Close()

		// Get all keys from bucket
		// TODO: Add the already existing names to the topic buckets (write a script ig?)
		aliasNames := []tui.FileData{}
		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte(topic))
			if bucket == nil {
				return fmt.Errorf("NO notes recorded yet...\n")
			}

			bucket.ForEach(func(k, v []byte) error {
				aliasNames = append(aliasNames, tui.FileData{
					Filename: string(k),
					Filepath: string(v),
				})
				return nil
			})
			return nil
		})

		if err != nil {
			return fmt.Errorf("KV error: %v\n", err)
		}

		// Sort based on filepath so all the aliases are together
		sort.Slice(aliasNames, func(i, j int) bool {
			return aliasNames[i].Filepath < aliasNames[j].Filepath
		})

		// TODO: Make this async
		fileContents := []string{}
		for _, file := range aliasNames {
			fileContent, err := os.ReadFile(file.Filepath)
			if err != nil {
				fileContents = append(fileContents, fmt.Sprintf("Couldn't read file: %v; Error: %v", file.Filename, err))
				continue
			}
			fileContents = append(fileContents, strings.TrimSpace(string(fileContent)))
		}

		// Start the bubbletea program to display the options
		reply := tui.FileData{
			Filename: "DIDN'T CHOOSE A FILE",
			Filepath: "DIDN'T CHOOSE A FILE",
		}
		model := tui.ListModelInitialize(aliasNames, fileContents, &reply)

		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Couldn't run bubbletea: %v\n", err)
			os.Exit(1)
		}
		if reply.Filename == "DIDN'T CHOOSE A FILE" && reply.Filepath == "DIDN'T CHOOSE A FILE" {
			return nil
		}

		s := tui.FileDisplay(reply)
		lipgloss.Println(s)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("topic", "t", "cyril", "Help message for toggle")
}
