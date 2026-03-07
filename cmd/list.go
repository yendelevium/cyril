/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	fp "path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/yendelevium/cyril/config"
	"github.com/yendelevium/cyril/internal/tui"
	bolt "go.etcd.io/bbolt"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all files (total and/or under a specific topic)",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		// Make this the TUI later
		for _, file := range aliasNames {
			log.Printf("%v %v", file.Filename, file.Filepath)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	listCmd.Flags().StringP("topic", "t", "cyril", "Help message for toggle")
}
