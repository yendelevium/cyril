/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cyril",
	Short: "A CLI butler to provide system-wide access to your notes. ",
	Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,

	// This allows the application to have both subcommands and arguments (loses the "did you mean command 'X' thing though...")
	// The subcommand takes precedence though so if you pass an argument with the same name as the subcommand ur goon
	Args: cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("Runs this when the rootCMD is executed %v", args)
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1) // Cobra supports persistent flags, which, if defined here,
		// will be global for your application.

		// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cyril.yaml)")
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here, will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cyril.yaml)")

	// Cobra also supports local flags, which will only run when this action is called directly.
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/cyril")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()
	viper.SetDefault("store", "$HOME/Documents/cyril")
	viper.SetDefault("editor", viper.Get("EDITOR"))

	var fileNotFoundError viper.ConfigFileNotFoundError
	err := viper.ReadInConfig()
	if err != nil {
		if errors.As(err, &fileNotFoundError) {
			log.Println("Using DEFAULT config (add your config at $HOME/.config/cyril/config.yml)")
		} else {
			cobra.CheckErr(err)
		}
	}
}
