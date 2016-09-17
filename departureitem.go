package main

import (
	"fmt"
	"time"

	"github.com/BenchR267/goalfred"
	"github.com/kiliankoe/dvbgo"
)

// wrapper for dvb.Departure
type departureItem dvb.Departure

// comply with goalfred's AlfredItem interface
func (dep departureItem) Item() goalfred.Item {
	var modeName string
	var modeTitle string
	mode, err := dvb.Departure(dep).Mode()
	if err != nil {
		modeName = "unknown"
		modeTitle = "Die"
	} else {
		modeName = mode.Name
		modeTitle = mode.Title
	}

	title := fmt.Sprintf("%s %s %s", dep.Line, dep.Direction, pluralizeTimeString(dep.RelativeTime))
	departureTime := time.Now().Add(time.Minute * time.Duration(dep.RelativeTime))

	notificationTime := time.Now().Add(time.Minute * time.Duration(dep.RelativeTime-notificationOffset))
	notificationDelay := notificationTime.Sub(time.Now()).Seconds()

	valid := true
	if notificationDelay < 0 {
		valid = false
	}

	item := goalfred.Item{
		Title:    title,
		Subtitle: formatSubtitleTime(departureTime),
		Valid:    &valid,
		Icon: &goalfred.Icon{
			Path: fmt.Sprintf("transport_icons/%s.png", modeName),
		},
	}

	item.SetComplexArg(goalfred.ComplexArg{
		Variables: map[string]interface{}{
			"notificationDelay": fmt.Sprintf("%.0f", notificationDelay),
			"line":              dep.Line,
			"direction":         dep.Direction,
			"modeTitle":         modeTitle,
			"departureTime":     departureTime.Format("15:04"),
			"notificationTime":  notificationTime.Format("15:04"),
		},
	})

	return item
}
