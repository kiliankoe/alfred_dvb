package main

import (
	"fmt"
	"os"

	"github.com/kiliankoe/dvbgo"
	"github.com/pascalw/go-alfred"
)

const (
	// Notifications are displayed 10 minutes before departure
	notificationOffset = 10
	// How many results are shown at a single time
	resultsAmount = 6
)

func main() {
	queryTerms := os.Args[1:]
	response := alfred.NewResponse()

	if len(queryTerms) < 1 {
		os.Exit(0)
	}

	departures, err := dvb.Monitor(queryTerms[0], 0, "")
	if err != nil {
		response.AddItem(&alfred.AlfredResponseItem{
			Title:    "Unerwarteter Fehler 😲",
			Subtitle: err.Error(),
		})
	} else if len(departures) < 1 {
		response.AddItem(&alfred.AlfredResponseItem{
			Title:    "Haltestelle nicht gefunden 🤔",
			Subtitle: "Vielleicht ein Tippfehler?",
		})
	} else {
		for _, dep := range departures[:resultsAmount] {
			mode, _ := dep.Mode()
			title := fmt.Sprintf("%s %s %s", dep.Line, dep.Direction, pluralizeTimeString(dep.RelativeTime))
			response.AddItem(&alfred.AlfredResponseItem{
				Title:    title,
				Subtitle: "",
				Arg:      "",
				Icon:     fmt.Sprintf("transport_icons/%s.png", mode.Name),
			})
		}
	}

	response.Print()
}

func pluralizeTimeString(minutes int) string {
	if minutes == 0 {
		return "jetzt"
	} else if minutes == 1 {
		return "in 1 Minute"
	}
	return fmt.Sprintf("in %d Minuten", minutes)
}
