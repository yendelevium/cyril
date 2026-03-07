/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	fp "path/filepath"

	"github.com/spf13/cobra"
	"github.com/yendelevium/cyril/config"
	bolt "go.etcd.io/bbolt"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new file/note under a topic",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the topic, else fallback to config.config.defaultTopic
		topic, err := cmd.Flags().GetString("topic")
		if err != nil {
			return err
		}
		if !cmd.Flags().Changed("topic") {
			topic = config.Conf.DefaultTopic
		}

		// Make all intermediate directories
		// dirpath := fmt.Sprintf("%s/%s", config.Conf.Store, topic)
		dirpath := fp.Join(config.Conf.Store, topic)
		// log.Println(dirpath)
		err = os.MkdirAll(dirpath, 0777)
		if err != nil {
			log.Fatalf("Couldn't create the required directories: %v", err)
		}

		// Create the file
		filename := args[0]
		filepath := fp.Join(dirpath, filename)

		// Check if it exists? If it does, return (we don't wanna override it)
		_, err = os.OpenFile(filepath, os.O_RDONLY, 0777)
		if err == nil {
			log.Printf("file:%s already exists in topic:%s; cannot create the file", filename, topic)
			// Returning the error makes the command help thingy pop up (idk)
			return nil
		}
		_, err = os.Create(filepath)
		if err != nil {
			log.Fatalf("Couldn't create file: %v", err)
		}

		// create the alias and store it in the DB
		done := make(chan struct{})
		go func() {
			AddAlias(filename, topic, filepath)
			// Signal the write is over
			done <- struct{}{}
		}()

		log.Printf("Created file; Topic: %s, Name: %s", topic, filename)
		// log.Println(config.config.editor)

		// TODO:  Hand-off control to the default editor... (maybe abstract this into its own function? probably)
		// The default editor might also not be there if $EDITOR isn't set, so add fallback to a universal editor (idk)
		command := exec.Command(config.Conf.Editor, filepath)

		// Need these coz otherwise the process starts elsewhere and NOT in the same terminal where cyril was invoked
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		// Run the command (editor) and wait for it to complete (Run() waits for compeletion automatically)
		if err = command.Run(); err != nil {
			return err
		}
		log.Println("Control returned")

		// Block till the DB write is done
		<-done
		return nil
	},
}

func init() {
	// Can't put config.config.defaultTopic as fallback value here
	// This is because I'm assuming this line executes before initconfig.Config() and we only get empty str (neither the default config.config, nor the actual config.config)
	// As config.config is initialized w/o any values..
	createCmd.Flags().StringP("topic", "t", "", "Help message for toggle")
	RootCmd.AddCommand(createCmd)
}

// TODO: Make this RETURN an error
func AddAlias(filename string, topic string, filepath string) {
	// Since BoltDB only allows one process to hold the file, we have to close the DB after every transaction...
	dbPath := fp.Join(config.Conf.DBPath, "cyril.db")
	db, err := bolt.Open(dbPath, 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("Couldn't open DB: %v", err)
	}
	defer db.Close()

	// Write alias name to the DB
	// The key is int the form {filename}.{topic}
	err = db.Update(func(tx *bolt.Tx) error {
		// Add in general
		bucket, err := tx.CreateBucketIfNotExists([]byte("cyril"))
		if err != nil {
			log.Fatalf("Failed to create bucket; %v", err)
		}

		bucketKey := fmt.Sprintf("%s.%s", filename, topic)
		err = bucket.Put([]byte(bucketKey), []byte(filepath))
		if err != nil {
			log.Fatalf("Failed to write to bucket; %v", err)
		}

		// Add for a specific topic
		topicBucket, err := tx.CreateBucketIfNotExists([]byte(topic))
		if err != nil {
			log.Fatalf("Failed to create bucket; %v", err)
		}
		err = topicBucket.Put([]byte(bucketKey), []byte(filepath))
		if err != nil {
			log.Fatalf("Failed to write to bucket; %v", err)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("KV error: %v", err)
	}

}
