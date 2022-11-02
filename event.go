package analytics

// Event is a product analytics event that is tracked.
type Event interface {
	userID() string
	// Type of the event, in the format 'namespace:entity.action.'
	//
	// namespace is the project namespace, like 'cf' for Common Fate.
	//
	// entity is the thing the event is related to, like 'scan'.
	//
	// action is the thing that happened in past tense, like 'created'.
	//
	// example type: "cf:scan.created"
	Type() string
	// Description of when the event is emitted.
	EmittedWhen() string
	// fixture generates a fixture event to be used in testing and examples.
	fixture()
}
