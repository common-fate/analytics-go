# analytics

This repository contains product analytics as defined in [RFD#8](https://github.com/common-fate/rfds/discussions/8).

## Events

All events emitted are listed in the below table. Discussion of these is welcomed on [RFD#8](https://github.com/common-fate/rfds/discussions/8).

| Name | Emitted When | Example Data |
| ---- | ----------- | ------------ |
| `cf:idp.synced` | IDP was synced | [./fixtures/cf-idp-synced.json](./fixtures/cf-idp-synced.json) |
| `cf:request.created` | Access Request was created | [./fixtures/cf-request-created.json](./fixtures/cf-request-created.json) |
| `cf:request.reviewed` | Access Request was reviewed | [./fixtures/cf-request-reviewed.json](./fixtures/cf-request-reviewed.json) |
| `cf:request.revoked` | Access Request was revoked | [./fixtures/cf-request-revoked.json](./fixtures/cf-request-revoked.json) |
| `cf:rule.archived` | Access Rule was archived | [./fixtures/cf-rule-archived.json](./fixtures/cf-rule-archived.json) |
| `cf:rule.created` | Access Rule was created | [./fixtures/cf-rule-created.json](./fixtures/cf-rule-created.json) |
| `cf:rule.updated` | Access Rule was updated | [./fixtures/cf-rule-updated.json](./fixtures/cf-rule-updated.json) |


## Usage

The `acore` package contains the core analytics client. The client is forked from the Rudderstack Go client (which itself appears to have been forked from Segment's Go SDK).

The library handles client-side hashing identifiers such as `usr_123`. We transform `usr_123` using a SHA256 hash into `usr_-CHh8_rdIqAotcBsP64GKQkfzW2hb1JDJ_6u7q4zom4` prior to events being dispatched. In the library this is controlled by the `analytics:"usr"` struct tag added to the event.

## Editing this README

This README is automatically generated from the template in [README.tmpl](./README.tmpl). Edit that rather than this file.
