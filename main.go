/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package main

import (
	"context"
	"os"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/fang"
	"github.com/yendelevium/cyril/cmd"
)

func main() {
	// If you don't wanna use charmbracelet/fang uncomment this line; also uncomment Execute in root.go
	// cmd.Execute()

	// I HATED the codeblock behind `usage`, and I changed that (along with the colour of `cyril` that comes under USAGE)
	var customColours fang.ColorSchemeFunc = func(lipgloss.LightDarkFunc) fang.ColorScheme {
		colorScheme := fang.DefaultColorScheme(lipgloss.LightDark(true))
		colorScheme.Codeblock = lipgloss.NoColor{}
		colorScheme.Program = lipgloss.RGBColor{
			R: 230,
			G: 185,
			B: 255,
		}
		return colorScheme
	}

	if err := fang.Execute(context.Background(), cmd.RootCmd, fang.WithoutVersion(), fang.WithoutCompletions(), fang.WithColorSchemeFunc(customColours)); err != nil {
		os.Exit(1)
	}
}
