/*
Copyright © 2026 yendelevium <yashbardia27@gmail.com>
*/
package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/yendelevium/cyril/cmd"
)

func main() {
	// If you don't wanna use charmbracelet/fang uncomment this line; also uncomment Execute in root.go
	// cmd.Execute()
	// TODO: I HATE the codeblock behind `usage`, change that somehow. HATE HATE HATEEE...
	if err := fang.Execute(context.Background(), cmd.RootCmd, fang.WithoutVersion(), fang.WithoutCompletions()); err != nil {
		os.Exit(1)
	}
}
