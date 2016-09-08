package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BenchR267/goalfred"
	"github.com/kiliankoe/dvbgo"
)

const (
	// Notifications are displayed 10 minutes before departure
	notificationOffset = 10
	// How many results are shown at a single time
	resultsAmount = 6
)

func main() {
	queryTerms := os.Args[1:]
	response := goalfred.NewResponse()

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

	// "albertplatz [3]" should filter everything but 3s
	var lineFilter = ""
	linefilterR, _ := regexp.Compile("\\[(.+)\\]")
	if matches := linefilterR.FindStringSubmatch(queryTerms[0]); len(matches) > 0 {
		stop = strings.Replace(stop, matches[0], "", -1)
		lineFilter = matches[1]
	}

	departures, err := dvb.Monitor(stop, offset, "")
	if lineFilter != "" {
		departures = filterDepartures(departures, lineFilter)
	}

	if err != nil {
		response.AddItem(&goalfred.Item{
			Title:    "Unerwarteter Fehler 😲",
			Subtitle: err.Error(),
		})
	} else if len(departures) < 1 {
		response.AddItem(&goalfred.Item{
			Title:    "Keine Haltestelle oder Verbindungen gefunden 🤔",
			Subtitle: "Vielleicht ein Tippfehler?",
		})
	} else {
		if len(departures) > resultsAmount {
			departures = departures[:resultsAmount]
		}

		for _, dep := range departures {
			response.AddItem(departureItem(*dep))
		}
	}

	response.Print()
}

type departureItem dvb.Departure

func (dep departureItem) Item() *goalfred.Item {
	var modeName string
	mode, err := dvb.Departure(dep).Mode()
	if err != nil {
		modeName = "unknown"
	} else {
		modeName = mode.Name
	}
	title := fmt.Sprintf("%s %s %s", dep.Line, dep.Direction, pluralizeTimeString(dep.RelativeTime))
	departureTime := time.Now().Add(time.Minute * time.Duration(dep.RelativeTime))

	notificationTime := time.Now().Add(time.Minute * time.Duration(dep.RelativeTime-notificationOffset))
	notificationDelay := notificationTime.Sub(time.Now()).Seconds()

	valid := true
	if notificationDelay < 0 {
		valid = false
	}

	item := &goalfred.Item{
		Title:    title,
		Subtitle: formatSubtitleTime(departureTime),
		Arg:      fmt.Sprintf("%.0f", notificationDelay),
		Valid:    &valid,
		Icon: &goalfred.Icon{
			Path: fmt.Sprintf("transport_icons/%s.png", modeName),
		},
	}
	return item
}

func filterDepartures(departures []*dvb.Departure, lineFilter string) []*dvb.Departure {
	var filtered []*dvb.Departure
	for _, dep := range departures {
		if strings.ToLower(dep.Line) == strings.ToLower(lineFilter) {
			filtered = append(filtered, dep)
		}
	}
	return filtered
}

func pluralizeTimeString(minutes int) string {
	if minutes == 0 {
		return "jetzt"
	} else if minutes == 1 {
		return "in 1 Minute"
	}
	return fmt.Sprintf("in %d Minuten", minutes)
}

func formatSubtitleTime(t time.Time) string {
	weekday := localizeWeekday(t.Weekday().String())
	minuteStr := ""
	if t.Minute() < 10 {
		minuteStr += "0"
	}
	minuteStr += strconv.Itoa(t.Minute())
	return fmt.Sprintf("%s, %d:%s Uhr", weekday, t.Hour(), minuteStr)
}

func localizeWeekday(weekday string) string {
	switch weekday {
	case "Monday":
		return "Montag"
	case "Tuesday":
		return "Dienstag"
	case "Wednesday":
		return "Mittwoch"
	case "Thursday":
		return "Donnerstag"
	case "Friday":
		return "Freitag"
	case "Saturday":
		return "Samstag"
	case "Sunday":
		return "Sonntag"
	default:
		return ""
	}
}
