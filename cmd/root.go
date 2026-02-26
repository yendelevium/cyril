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

type Config struct {
	Store        string `mapstructure:"store"`
	Editor       string `mapstructure:"editor"`
	DefaultTopic string `mapstructure:"defaultTopic"`
	DBPath       string `mapstructure:"dbPath"`
}

var Conf Config
var hey string

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
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
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/cyril")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	var fileNotFoundError viper.ConfigFileNotFoundError
	err := viper.ReadInConfig()
	if err != nil {
		if errors.As(err, &fileNotFoundError) {
			log.Println("Using DEFAULT config (add your config at $HOME/.config/cyril/config.yml)")
		} else {
			cobra.CheckErr(err)
		}
	}

	// Setting default configs; viper.SetDefault doesn't work with viper.Unmarshall -> https://github.com/spf13/viper/issues/1284
	// viper.SetDefault("store", "$HOME/Documents/cyril")
	// viper.SetDefault("editor", viper.Get("EDITOR"))

	Conf.Store = `$HOME/Documents/cyrilStore`
	Conf.Editor = `$EDITOR`
	Conf.DefaultTopic = "general"
	Conf.DBPath = Conf.Store // store the db at the same place as the notes by default

	err = viper.Unmarshal(&Conf)
	if err != nil {
		log.Println("Unmarshal failed?")
	}

	// Make sure the values are replaced by the environment values
	Conf.Store = os.ExpandEnv(Conf.Store)
	Conf.Editor = os.ExpandEnv(Conf.Editor)
	Conf.DBPath = os.ExpandEnv(Conf.DBPath)
}
