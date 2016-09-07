package main

import (
	"fmt"
	"os"

	"github.com/kiliankoe/dvbgo"
	"github.com/pascalw/go-alfred"
)

func main() {
	queryTerms := os.Args[1:]
	response := alfred.NewResponse()

	if len(queryTerms) < 1 {
		os.Exit(0)
	}

	departures, err := dvb.Monitor(queryTerms[0], 0, "")
	if err != nil {
		fmt.Println(err) // FIXME
	}

	if len(departures) < 1 {
		response.AddItem(&alfred.AlfredResponseItem{
			Valid: true,
			Uid:   "",
			Title: "Haltestelle nicht gefunden",
		})
	} else {
		for _, dep := range departures[:6] {
			response.AddItem(&alfred.AlfredResponseItem{
				Valid: true,
				Uid:   dep.String(),
				Title: dep.String(),
				Arg:   "",
			})
		}
	}

	response.Print()
}
