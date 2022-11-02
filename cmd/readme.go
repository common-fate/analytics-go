// Program readme generates the table of contents in the
// README for this package
package main

import (
	"html/template"
	"log"
	"os"

	"github.com/common-fate/analytics-go"
	"github.com/common-fate/analytics-go/internal"
)

type event struct {
	Name        string
	EmittedWhen string
	FixturePath string
}

// context is used as template data
type context struct {
	Events []event
}

func main() {

	t, err := template.ParseFiles("README.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	var events []event
	for _, e := range analytics.AllEvents {
		t := e.Type()
		evt := event{
			Name:        t,
			FixturePath: "./" + internal.FixturePath(t),
			EmittedWhen: e.EmittedWhen(),
		}
		events = append(events, evt)
	}

	f, err := os.OpenFile("README.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = t.Execute(f, context{Events: events})
	if err != nil {
		log.Fatal(err)
	}
}
