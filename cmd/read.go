/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read a saved note, either by its name or its alias",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]
		aliasNames := make(map[string]string)

		// Get all keys that have the same alias name but diff topics
		matchAliasPrefixes(filename, aliasNames)

		// TODO: Walk the store to see any matching files using os.Walkdir (I think)
		// TODO: IF file doesn't exist offer to create it
		// TODO: IF multiple files exist offer to choose which one to open
		if len(aliasNames) == 0 {
			// No alias names
			// Walk the directory? or NO?
			log.Printf("Couldn't find the note %s", filename)
			return nil
		}

		// Iterate through all aliasnames and print out their content -> make the user choose in the future
		for key, value := range aliasNames {
			// log.Println(value)
			fileContent, err := os.ReadFile(value)
			if err != nil {
				log.Fatalf("Couldn't read file: %v; Error: %v", key, err)
			}
			log.Printf("%s", string(fileContent))

			// For now only reading the first aliased file
			break
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(readCmd)
}

func matchAliasPrefixes(filename string, aliasNames map[string]string) {
	// Making this a function as I want to use defer statement so I don't forget to close the DB
	// Also using it in ReadOnly mode
	dbPath := fmt.Sprintf("%s/cyril.db", Conf.DBPath)
	db, err := bolt.Open(dbPath, 0644, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	if err != nil {
		log.Fatalf("Couldn't open DB: %v", err)
	}
	defer db.Close()

	// Get all keys with the filename as an alias
	// Stolen from official docs
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("cyril"))
		if bucket == nil {
			return errors.New("NO notes recorded yet...")
		}

		// Matching Prefixes (multiple topics might have same filename)
		c := bucket.Cursor()
		prefix := []byte(filename)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			aliasNames[string(k)] = string(v)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("KV error: %v", err)
	}
}
