package resources

import (
	"glow-gui/glow"
	"strings"
)

type LabelID uint16

const (
	GlowLabel LabelID = iota
	ColumnsLabel
	LengthLabel
	RowsLabel
	IntervalLabel
	LayersLabel
	GridLabel
	ColorsLabel
	HueShiftLabel
	ScanLengthLabel
	BeginLabel
	EndLabel
	OriginLabel
	OrientationLabel
	HueLabel
	ScanLabel
	SaturationLabel
	ValueLabel
	EffectsLabel
	DynamicLabel
	StaticLabel
	GradientLabel
	RevertLabel
	ApplyLabel
	RateLabel
	OverrideLabel
	CutLabel
	CopyLabel
	PasteLabel
	EditLabel
	PickerLabel
)

var entryLabels = []string{
	"Glow", "Columns",
	"Length", "Rows", "Interval", "Layers",
	"Grid", "Colors", "Hue Shift", "Scan Length",
	"Begin", "End",
	"Origin", "Orientation",
	"Hue", "Scan",
	"Saturation", "Value",
	"Effects",
	"Dynamic", "Static", "Gradient",
	"Revert", "Apply",
	"Interval (ms)", "Override",
	"Cut", "Copy", "Paste", "Edit", "Picker",
}

func (id LabelID) String() string {
	return entryLabels[id]
}

func (id LabelID) PlaceHolder() string {
	return strings.ToLower(entryLabels[id])
}

type OrientationID glow.Orientation

var OrientationLabels = []string{
	"Level",
	"Upright",
	"Tilted",
}

func (id OrientationID) String() string {
	return OrientationLabels[id]
}

func (id OrientationID) PlaceHolder() string {
	return strings.ToLower(OrientationLabels[id])
}

type OriginID glow.Origin

var OriginLabels = []string{
	"Top Left",
	"Top Right",
	"Bottom Left",
	"Bottom Right",
}

func (id OriginID) String() string {
	return OriginLabels[id]
}

func (id OriginID) PlaceHolder() string {
	return strings.ToLower(OriginLabels[id])
}