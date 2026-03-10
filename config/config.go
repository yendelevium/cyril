package config

import (
	"errors"
	"log"
	"os"
	"os/exec"

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

func InitConfig() {
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

	// Potential fallback editors if not found in config or if $EDITOR isn't set
	if Conf.Editor == "" {
		editors := []string{"vim", "emacs", "nano"}
		for _, name := range editors {
			// If it exists, set it and retrun
			path, err := exec.LookPath(name)
			if err == nil {
				Conf.Editor = path
				return
			}
		}

		// If no editors found, exit process
		log.Fatalf("Failed to find a usable editor, please set $EDITOR or an editor in ~/.config/cyril/config.yml")
	}
}
