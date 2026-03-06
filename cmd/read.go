/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"time"

	fp "path/filepath"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
	"github.com/spf13/cobra"
	"github.com/yendelevium/cyril/config"
	"github.com/yendelevium/cyril/internal/tui"
	bolt "go.etcd.io/bbolt"
)

// TODO: Make this common between the tui and this?
type fileData struct {
	filename string
	filepath string
}

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read a saved note, either by its name or its alias",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]
		aliasNames := []fileData{}

		// Get all keys that have the same alias name but diff topics
		err := MatchAliasPrefixes(filename, &aliasNames)
		if err != nil {
			return err
		}

		// TODO: Walk the store to see any matching files using os.Walkdir (I think)
		// TODO: IF file doesn't exist offer to create it
		// TODO: IF multiple files exist offer to choose which one to open
		// If aliasName length is 0 or 1, handle normally
		if len(aliasNames) == 0 {
			// No alias names
			// Walk the directory? or NO?
			fmt.Printf("Couldn't find the note %s!\n", filename)
			return nil
		}

		// TODO: This lipgloss code is reused in the tui View() as well, so un-reuse it by making it a function?
		if len(aliasNames) == 1 {
			file := aliasNames[0]
			fmt.Printf("Alias: %s; Path: %s\n", file.filename, file.filepath)
			fileContent, err := os.ReadFile(file.filepath)
			if err != nil {
				fmt.Printf("Couldn't read file: %v; Error: %v\n", file.filename, err)
				return nil
			}

			// Get the current terminal width and use that to help deal with border rendering issues in lipgloss
			// TODO: This fix works, but then if I fullscreen the terminal, I get a block that is half rendered and I can't do anything coz it has alr been printed and the program has ended
			// It isn't a HUGE problem but still it looks kindof ugly so yeah idk
			physicalWidth, _, err := term.GetSize(os.Stdout.Fd())
			if err != nil {
				physicalWidth = 80 // Fallback width just in case
			}
			style := lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				// TODO: Factor the colour out as a constant as I'm using it in multiple places
				BorderForeground(lipgloss.RGBColor{
					R: 220,
					G: 155,
					B: 255,
				}).
				Width(physicalWidth - 1).
				PaddingLeft(1).
				PaddingRight(1)
			lipgloss.Println(style.Render(fmt.Sprintf("%s", string(fileContent))))
			return nil
		}

		// Iterate through all aliasnames and print out their content -> make the user choose
		model := tui.ReadModel{}
		for _, file := range aliasNames {
			// fmt.Printf("Idx: %d; Alias: %s; Path: %s", idx, file.filename, file.filepath)
			model.Files = append(model.Files, tui.FileData{
				Filename: file.filename,
				Filepath: file.filepath,
			})
		}

		// Start the bubbletea program to display the options
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Couldn't run bubbletea: %v\n", err)
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(readCmd)
}

func MatchAliasPrefixes(filename string, aliasNames *[]fileData) error {
	// Making this a function as I want to use defer statement so I don't forget to close the DB
	// Also using it in ReadOnly mode
	// dbPath := fmt.Sprintf("%s/cyril.db", Conf.DBPath)
	dbPath := fp.Join(config.Conf.DBPath, "cyril.db")
	db, err := bolt.Open(dbPath, 0644, &bolt.Options{Timeout: 1 * time.Second, ReadOnly: true})
	if err != nil {
		return fmt.Errorf("Couldn't open DB: %v \n", err)
	}
	defer db.Close()

	// Get all keys with the filename as an alias
	// Stolen from official docs
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("cyril"))
		if bucket == nil {
			return fmt.Errorf("NO notes recorded yet...\n")
		}

		// Matching Prefixes (multiple topics might have same filename)
		c := bucket.Cursor()
		prefix := []byte(filename)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			*aliasNames = append(*aliasNames, fileData{
				filename: string(k),
				filepath: string(v),
			})
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("KV error: %v\n", err)
	}
	return nil
}
