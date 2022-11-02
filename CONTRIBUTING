# Contributing

## TLDR

To add a new event:

1. Copy/paste an existing `event_*.go` file.
2. Edit the event, following the patterns in the file.
3. Update the test fixtures with `go test -update`.
4. Run tests with `go test ./...`.
5. Update the README with `go run cmd/readme.go`.

## Adding Events (the long explanation)

To add a new event, add a new `.go` file in the root of the repo named after your event in `snake_case` format. For example, if we'd like to add a `cf:scan.created` event, create a file called `event_scan_created.go`. This file will contain a struct for your event.

Your file must call `registerEvent` in an init function as shown below.

```go
// event_scan_created.go

func init() {
	registerEvent(&ScanCreated{})
}

type ScanCreated struct {}
```

Your event must implement the `Event` interface:

```go
// Event is a product analytics event that is tracked.
type Event interface {
	userID() string
	Type() string
	EmittedWhen() string
	fixture()
}
```

The `userID` method should return the ID of the user that created the event. The `Type` should return the name of the event.

We use the following event name format:

```
namespace:entity.action
```

Where:

- `namespace` is the project namespace, like `cf` for the Common Fate repository
- `entity` is the thing the event relates to, like `scan`
- `action` is the thing that happened, like `created`

Additionally, you **must** add an `analytics` struct tag onto any fields which contain identifiers. This implements client-side hashing to avoid sending any raw identifiers which correspond to database values.

You also must add `json` tags to your struct fields using `snake_case` format.

An example for the `cf:scan.created` event might look like this:

```go
// event_scan_created.go

func init() {
	registerEvent(&ScanCreated{})
}

type ScanCreated struct {
	CreatedBy string     `json:"created_by" analytics:"usr"`
}

func (e *ScanCreated) userID() string { return e.CreatedBy }

func (e *ScanCreated) Type() string { return "cf:scan.created" }

func (e *ScanCreated) EmittedWhen() string { return "Scan was created" }

func (e *ScanCreated) fixture() {
	*e = ScanCreated{
		CreatedBy: "usr_123"
	}
}
```

## Generating fixtures

This package uses snapshot testing to ensure events are consistent and to explain what data is being sent in the anonymous analytics.

After adding a new event, run the following command to update fixture data:

```
go test -update
```

Make sure to inspect your fixture data after it's created to ensure it's in line with what you expect, as this is what will be dispatched.

## Updating the README

The README is automatically generated from the available events. To update it after adding a new event, run:

```
go run cmd/readme.go
```
