/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package cmd

import (
	"io"
	"log"
	"os"
	"os/exec"
	fp "path/filepath"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the target note",
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
		MatchAliasPrefixes(filename, aliasNames)

		// TODO: The same ones as in read.go (walk file dir, choose if multiple, prompt to create)
		if len(aliasNames) == 0 {
			log.Printf("Couldn't find the note %s; do you want to create it? {make it a prompt}", filename)
			return nil
		}

		// TODO: Only ONE file must be edited haha
		for key, value := range aliasNames {
			err := editFile(key, value)
			if err != nil {
				return err
			}
			// For now only editing the first aliased file
			break
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}

// Copy a file from src to dst
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return nil
}

// Facing some concurrency issues
// IF another process opens the same file and starts editing, I'll be in trouble
// Because both will create the same file and when one is done, it will delete that temp file, but the other process might still be writing to it
// I can't do anything if 2 things edit the same file, but same alias I can
// I'll name the file by its alias name {filename}.{topicname} so they won't create the same temp file and can safely edit each other
// TODO: Store all tmp files of cyril under /tmp/cyril ig?
func editFile(filename string, filepath string) error {
	// Create the temporary directory and ensure we clean up when done.
	tmpDir := os.TempDir()

	tmpFile := fp.Join(tmpDir, filename)
	if err := copyFile(filepath, tmpFile); err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	// TODO:  Hand-off control to the default editor... (maybe abstract this into its own function? probably)
	// The default editor might also not be there if $EDITOR isn't set, so add fallback to a universal editor (idk)
	command := exec.Command(Conf.Editor, tmpFile)

	// Need these coz otherwise the process starts elsewhere and NOT in the same terminal where cyril was invoked
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Run the command (editor) and wait for it to complete (Run() waits for compeletion automatically)
	if err := command.Run(); err != nil {
		return err
	}
	log.Println("Control returned")

	// Copy the tmp file back to the original file
	return copyFile(tmpFile, filepath)
}
