package analytics

var AllEvents = map[string]Event{}

func registerEvent(e Event) {
	AllEvents[e.Type()] = e
}
