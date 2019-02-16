package component

import "encoding/json"

// FlexLayoutItem is an item in a flex layout.
type FlexLayoutItem struct {
	Width int           `json:"width,omitempty"`
	View  ViewComponent `json:"view,omitempty"`
}

// FlexLayoutSection is a slice of items group together.
type FlexLayoutSection []FlexLayoutItem

// FlexLayoutConfig is configuration for the flex layout view.
type FlexLayoutConfig struct {
	Sections []FlexLayoutSection `json:"sections,omitempty"`
}

// FlexLayout is a flex layout view.
type FlexLayout struct {
	base
	Config FlexLayoutConfig `json:"config,omitempty"`
}

// NewFlexLayout creates an instance of FlexLayout.
func NewFlexLayout(title string) *FlexLayout {
	return &FlexLayout{
		base: newBase(typeFlexLayout, TitleFromString(title)),
	}
}

// GetMetdata returns the metadata for the flex layout view.
func (fl *FlexLayout) GetMetadata() Metadata {
	return fl.Metadata
}

// AddSections adds one or more sections to the flex layout.
func (fl *FlexLayout) AddSections(sections ...FlexLayoutSection) {
	fl.Config.Sections = append(fl.Config.Sections, sections...)
}

type flexLayoutMarshal FlexLayout

// MarshalJSON marshals the flex layout to JSON.
func (fl *FlexLayout) MarshalJSON() ([]byte, error) {
	x := flexLayoutMarshal(*fl)
	x.Metadata.Type = typeFlexLayout
	return json.Marshal(&x)
}
