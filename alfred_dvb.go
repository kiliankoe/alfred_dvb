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

// Notifications are displayed 10 minutes before departure or whatever is specified in Alfred
func getNotificationOffset() int {
	offset := os.Getenv("TIME_OFFSET")
	intVal, err := strconv.Atoi(offset)
	if err != nil {
		return 10
	}
	return intVal
}

// Default is to show 6 results, unless otherwise specified in Alfred
func getResultsAmount() int {
	amount := os.Getenv("RESULTS_AMOUNT")
	intVal, err := strconv.Atoi(amount)
	if err != nil {
		return 6
	}
	return intVal
}

var (
	notificationOffset = getNotificationOffset()
	resultsAmount      = getResultsAmount()
)

func main() {
	queryTerms := os.Args[1:]
	if len(queryTerms) < 1 {
		os.Exit(0)
	}

	query, err := goalfred.Normalize(queryTerms[0])
	stop, offset, lineFilter := parseQuery(query)

	departures, err := dvb.Monitor(stop, offset, "")
	if lineFilter != "" {
		departures = filterDepartures(departures, lineFilter)
	}

	response := goalfred.NewResponse()

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
		Valid:    &valid,
		Icon: &goalfred.Icon{
			Path: fmt.Sprintf("transport_icons/%s.png", modeName),
		},
	}

	item.SetComplexArg(goalfred.ComplexArg{
		Variables: map[string]string{
			"notificationDelay": fmt.Sprintf("%.0f", notificationDelay),
			"line":              dep.Line,
			"direction":         dep.Direction,
			"departureTime":     fmt.Sprintf("%02d:%02d Uhr", departureTime.Hour(), departureTime.Minute()),
			"notificationTime":  fmt.Sprintf("%02d:%02d Uhr", notificationTime.Hour(), notificationTime.Minute()),
		},
	})

	return item
}

func parseQuery(query string) (stop string, offset int, lineFilter string) {
	stop = query

	// "helmholtzstraße in 10" should set the offset accordingly
	offsetR, _ := regexp.Compile("in (\\d+)")
	if matches := offsetR.FindStringSubmatch(query); len(matches) > 0 {
		stop = strings.Split(stop, " in ")[0]
		offset, _ = strconv.Atoi(matches[1])
	}

	// "albertplatz [3]" should filter everything but 3s
	linefilterR, _ := regexp.Compile("\\[(.+)\\]")
	if matches := linefilterR.FindStringSubmatch(query); len(matches) > 0 {
		stop = strings.Replace(stop, matches[0], "", -1)
		lineFilter = matches[1]
	}

	return
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
	return fmt.Sprintf("%s, %02d:%02d Uhr", weekday, t.Hour(), t.Minute())
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
