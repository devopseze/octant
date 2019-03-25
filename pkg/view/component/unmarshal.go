package component

import (
	"encoding/json"

	"github.com/pkg/errors"
)

func unmarshal(to TypedObject) (interface{}, error) {
	var o interface{}
	var err error

	switch to.Metadata.Type {
	case typeContainers:
		t := &Containers{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal containers config")
		o = t
	case typeExpressionSelector:
		t := &ExpressionSelector{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal expressionSelector config")
		o = t
	case typeFlexLayout:
		t := &FlexLayout{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal expressionSelector config")
		o = t
	case "grid":
		t := &Grid{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal grid config")
		o = t
	case typeLabels:
		t := &Labels{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal labels config")
		o = t
	case typeLabelSelector:
		t := &LabelSelector{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal labelSelector config")
		o = t
	case typeLink:
		t := &Link{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal link config")
		o = t
	case typeList:
		t := &List{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal list config")
		o = t
	case "panel":
		t := &Panel{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal panel config")
		o = t
	case typeQuadrant:
		t := &Quadrant{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal quadrant config")
		o = t
	case typeResourceViewer:
		t := &ResourceViewer{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal resourceViewer config")
		o = t
	case typeSelectors:
		t := &Selectors{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal selectors config")
		o = t
	case typeSummary:
		t := &Summary{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal summary config")
		o = t
	case typeTable:
		t := &Table{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal table config")
		o = t
	case typeText:
		t := &Text{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal text config")
		o = t
	case typeTimestamp:
		t := &Timestamp{base: base{Metadata: to.Metadata}}
		err = errors.Wrapf(json.Unmarshal(to.Config, &t.Config),
			"unmarshal timestamp config")
		o = t

	default:
		return nil, errors.Errorf("unknown view component %q", to.Metadata.Type)
	}

	if err != nil {
		return nil, errors.Wrapf(err, "unmarshal component")
	}

	return o, nil
}
