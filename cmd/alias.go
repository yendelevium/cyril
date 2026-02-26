/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

// aliasCmd represents the alias command
var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Create an alias for a file",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),

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
		aliasNames := make(map[string]string)
		MatchAliasPrefixes(targetAlias, aliasNames)
		if len(aliasNames) == 0 {
			// No alias names
			// Walk the directory? or NO?
			log.Printf("Couldn't find the note %s", targetAlias)
			return nil
		}

		// Iterate through all aliasnames and steal the topic -> make the user choose in the future
		for key, value := range aliasNames {
			log.Printf("Alias: %s; Path: %s", key, value)
			keyParts := strings.Split(key, ".")

			// TODO: This will FAIL if the topic has a '.' in it, so uhh idk fix this later?
			// Maybe I can force topics to NOT contain any '.'s
			topic := keyParts[len(keyParts)-1]

			// Add the alias
			AddAlias(newAlias, topic, value)
			log.Printf("New alias for file: %s in topic: %s has been created", targetAlias, topic)

			// For now only reading the first alias
			break
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(aliasCmd)
}
