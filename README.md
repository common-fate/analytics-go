# analytics

This repository contains product analytics as defined in [RFD#8](https://github.com/common-fate/rfds/discussions/8).

## Event types

The [fixtures](./fixtures/) folder contains examples of each event type which is dispatched, along with the properties.

## Usage

The `acore` package contains the core analytics client. The client is forked from the Rudderstack Go client (which itself appears to have been forked from Segment's Go SDK).

The library handles client-side hashing identifiers such as `usr_123`. We transform `usr_123` using a SHA256 hash into `usr_-CHh8_rdIqAotcBsP64GKQkfzW2hb1JDJ_6u7q4zom4` prior to events being dispatched. In the library this is controlled by the `analytics:"usr"` struct tag added to the event.
