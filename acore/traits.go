package acore

// This type is used to represent traits in messages that support it.
// It is a free-form object so the application can set any value it sees fit but
// a few helper method are defined to make it easier to instantiate traits with
// common fields.
// Here's a quick example of how this type is meant to be used:
//
//	analytics.Identify{
//		UserId: "0123456789",
//		Traits: analytics.NewTraits()
//			.Set("Role", "Jedi"),
//	}
//

type Traits map[string]interface{}

func NewTraits() Traits {
	return make(Traits, 10)
}

func (t Traits) Set(field string, value interface{}) Traits {
	t[field] = value
	return t
}
