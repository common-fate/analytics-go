package analytics

type Provider struct {
	Publisher string `json:"publisher"`
	Name      string `json:"name"`
	Version   string `json:"version"`

	// Kind is only used for target groups
	Kind string `json:"kind,omitempty"`
}
