package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

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

	stop := queryTerms[0]
	offset := 0

	// "helmholtzstraße in 10" should set the offset accordingly
	offsetR, _ := regexp.Compile("in (\\d+)")
	if matches := offsetR.FindStringSubmatch(queryTerms[0]); len(matches) > 0 {
		stop = strings.Split(stop, " in ")[0]
		offset, _ = strconv.Atoi(matches[1])
	}

	departures, err := dvb.Monitor(stop, offset, "")
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
