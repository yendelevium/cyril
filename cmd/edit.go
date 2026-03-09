/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	fp "path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
	"github.com/yendelevium/cyril/config"
	"github.com/yendelevium/cyril/internal/tui"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "edit the target note",
	Long:  "Lets you edit the given note. You can access it via the filename or the aliasname. If there are multiple files with the same name across topics/multiple prefix-matches for that file, you can choose which one to edit via a selector.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]
		aliasNames := []tui.FileData{}

		// Get all keys that have the same alias name but diff topics
		err := MatchAliasPrefixes(filename, &aliasNames)
		if err != nil {
			return err
		}

		// TODO: The same ones as in read.go (walk file dir, choose if multiple, prompt to create)
		if len(aliasNames) == 0 {
			log.Printf("Couldn't find the note %s; do you want to create it? {make it a prompt}", filename)
			return nil
		}

		// IF only one file, edit & return
		if len(aliasNames) == 1 {
			file := aliasNames[0]
			err := editFile(file.Filename, file.Filepath)
			if err != nil {
				return err
			}
			return nil
		}

		// Iterate through all aliasnames and print out their content -> make the user choose
		model := tui.EditModel{
			Files: aliasNames,
			Reply: &tui.FileData{
				Filename: "DIDN'T CHOOSE A FILE",
				Filepath: "DIDN'T CHOOSE A FILE",
			},
			NoOption: false,
		}

		// Start the bubbletea program to display the options
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Couldn't run bubbletea: %v\n", err)
			os.Exit(1)
		}

		if model.Reply.Filename == "DIDN'T CHOOSE A FILE" && model.Reply.Filepath == "DIDN'T CHOOSE A FILE" {
			return nil
		}

		// Print it out
		err = editFile(model.Reply.Filename, model.Reply.Filepath)
		if err != nil {
			return err
		}

		fmt.Println("Successfully edited the file!")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(editCmd)
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
	command := exec.Command(config.Conf.Editor, tmpFile)

	// Need these coz otherwise the process starts elsewhere and NOT in the same terminal where cyril was invoked
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// Run the command (editor) and wait for it to complete (Run() waits for compeletion automatically)
	if err := command.Run(); err != nil {
		return err
	}
	// Copy the tmp file back to the original file
	return copyFile(tmpFile, filepath)
}
