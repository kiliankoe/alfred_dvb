package main

import (
	"os"

	"github.com/BenchR267/goalfred"
)

// If an error is given, output that and stop all further execution
func handleError(err error) {
	if err == nil {
		return
	}

	goalfred.Add(goalfred.Item{
		Title:    "Unerwarteter Fehler 😲",
		Subtitle: err.Error(),
	})
	goalfred.Print()

	os.Exit(1)
}
