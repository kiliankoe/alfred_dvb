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

var (
	notificationOffset = getNotificationOffset()
	resultsAmount      = getResultsAmount()
)

func main() {
	queryTerms := os.Args[1:]
	if len(queryTerms) < 1 {
		// Only continue with further execution if an argument is given
		os.Exit(0)
	}

	// Re-normalize input to deal with Alfred issues
	// See http://www.alfredforum.com/topic/2015-encoding-issue/ for details
	query, err := goalfred.Normalize(queryTerms[0])
	handleError(err)

	stop, offset, lineFilter := parseQuery(query)

	departures, err := dvb.Monitor(stop, offset, "")
	handleError(err)

	if lineFilter != "" {
		departures = filterDepartures(departures, lineFilter)
	}

	defer goalfred.Print()

	if len(departures) < 1 {
		goalfred.Add(goalfred.Item{
			Title:    "Keine Haltestelle oder Verbindungen gefunden 🤔",
			Subtitle: "Vielleicht ein Tippfehler?",
		})
		return
	}

	if len(departures) > resultsAmount {
		departures = departures[:resultsAmount]
	}

	for _, dep := range departures {
		goalfred.Add(departureItem(*dep))
	}
}

// Notifications are displayed 10 minutes before departure unless otherwise specified in Alfred
func getNotificationOffset() int {
	offset := os.Getenv("TIME_OFFSET")
	intVal, err := strconv.Atoi(offset)
	if err != nil {
		return 10
	}
	return intVal
}

// Default is to show 6 results unless otherwise specified in Alfred
func getResultsAmount() int {
	amount := os.Getenv("RESULTS_AMOUNT")
	intVal, err := strconv.Atoi(amount)
	if err != nil {
		return 6
	}
	return intVal
}

// Read name of the stop, optional time offset and optional line filter from query
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

// Filter given list of departures by given line
func filterDepartures(departures []*dvb.Departure, lineFilter string) []*dvb.Departure {
	var filtered []*dvb.Departure
	for _, dep := range departures {
		if strings.ToLower(dep.Line) == strings.ToLower(lineFilter) {
			filtered = append(filtered, dep)
		}
	}
	return filtered
}

// Given a number of minutes, output the correct version of "in x minute(s)" or "now"
func pluralizeTimeString(minutes int) string {
	if minutes == 0 {
		return "jetzt"
	} else if minutes == 1 {
		return "in 1 Minute"
	}
	return fmt.Sprintf("in %d Minuten", minutes)
}

// Given a time object, output something like "Monday 15:04 Uhr"
func formatSubtitleTime(t time.Time) string {
	weekday := localizeWeekday(t.Weekday().String())
	return fmt.Sprintf("%s %s Uhr", weekday, t.Format("15:04"))
}

// Translate a weekday string to German
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
