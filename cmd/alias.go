/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yendelevium/cyril/internal/tui"
)

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "create an alias for an existing file",
	Long:  "Adds a new alias name for the given file (or an existing aliasname) in the same topic where the file belongs.",
	Args:  cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		// Read the alias value from the KV
		// Get the topic from the alias name, create the key to store
		// Make the new alias
		// TODO: Warn if there is an existing alias of the same name pointing to a diff file?
		// TODO: If multiple with same alias, prompt the user which one he wants
		// TODO: Add -t tag
		// `cyril alias <new_alias> <old_alias>`
		newAlias := args[0]
		targetAlias := args[1]

		// Get the other alias(es)
		aliasNames := []tui.FileData{}
		MatchAliasPrefixes(targetAlias, &aliasNames, "cyril")
		if len(aliasNames) == 0 {
			// No alias names
			// Walk the directory? or NO?
			log.Printf("Couldn't find the note %s", targetAlias)
			return nil
		}

		// Iterate through all aliasnames and steal the topic -> make the user choose in the future
		for _, file := range aliasNames {
			log.Printf("Alias: %s; Path: %s", file.Filename, file.Filepath)
			keyParts := strings.Split(file.Filename, ".")

			// TODO: This will FAIL if the topic has a '.' in it, so uhh idk fix this later?
			// Maybe I can force topics to NOT contain any '.'s
			topic := keyParts[len(keyParts)-1]

			// Add the alias
			AddAlias(newAlias, topic, file.Filepath)
			log.Printf("New alias for file: %s in topic: %s has been created", targetAlias, topic)

			// For now only reading the first alias
			break
		}
		return nil
	},
}

func init() {
	RootCmd.AddCommand(aliasCmd)
}
