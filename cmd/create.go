/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new note",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Get the topic, else fallback to config.defaultTopic
		topic, err := cmd.Flags().GetString("topic")
		if err != nil {
			return err
		}
		if !cmd.Flags().Changed("topic") {
			topic = config.defaultTopic
		}

		// Make all intermediate directories
		dirpath := fmt.Sprintf("%s/%s", config.store, topic)
		// log.Println(dirpath)
		err = os.MkdirAll(dirpath, 0777)
		if err != nil {
			log.Fatalf("Couldn't create the required directories: %v", err)
		}

		// Create the file
		filename := args[0]
		filepath := fmt.Sprintf("%s/%s", dirpath, filename)
		_, err = os.Create(filepath)
		if err != nil {
			log.Fatalf("Couldn't create file: %v", err)
		}

		log.Printf("Created file; Topic: %s, Name: %s", topic, filename)
		return nil
	},
}

func init() {
	// Can't put config.defaultTopic as fallback value here
	// This is because I'm assuming this line executes before initConfig() and we only get empty str (neither the default config, nor the actual config)
	// As config is initialized w/o any values..
	createCmd.Flags().StringP("topic", "t", "", "Help message for toggle")
	rootCmd.AddCommand(createCmd)
}
