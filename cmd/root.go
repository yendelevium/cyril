/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/yendelevium/cyril/config"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cyril",
	Short: "A CLI butler to provide system-wide access to your notes. ",
	Long:  "The CLI butler you didn't know you needed. Get system-wide access to your notes with a minimalistic TUI :) Create, read, edit and give aliases to your notes so you never have to dig through hundreds of notes trying to find the one you need. It'll be right where you work - in your terminal.",

	// This allows the application to have both subcommands and arguments (loses the "did you mean command 'X' thing though...")
	// The subcommand takes precedence though so if you pass an argument with the same name as the subcommand ur goon
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Runs this when the rootCMD is executed %v", args)
		return nil
	},
}

// No need fr this as I'm now using charmbracelet/fang, but if you want to revert just uncomment; also uncomment main.go
// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// func Execute() {
// 	err := RootCmd.Execute()
// 	if err != nil {
// 		os.Exit(1)
// 	}
// }

func init() {
	cobra.OnInitialize(config.InitConfig)
}
