package ui

import (
	"glow-gui/data"
	"glow-gui/glow"
	"glow-gui/resources"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FrameEditor struct {
	*fyne.Container
	model       *data.Model
	frame       *glow.Frame
	layerSelect *widget.Select
	fields      *data.FrameFields
	rateBounds  *IntEntryBounds
	rateBox     *RangeIntBox
	tools       *FrameTools
}

func NewFrameEditor(model *data.Model, window fyne.Window,
	sharedTools *SharedTools) *FrameEditor {

	fe := &FrameEditor{
		model:       model,
		layerSelect: NewLayerSelect(model),
		rateBounds:  RateBounds,
		fields:      data.NewFrameFields(),
		frame:       &glow.Frame{},
	}

	fe.layerSelect = NewLayerSelect(fe.model)
	ratelabel := widget.NewLabel(resources.RateLabel.String())
	fe.rateBox = NewRangeIntBox(fe.fields.Interval, fe.rateBounds)
	fe.fields.Interval.AddListener(binding.NewDataListener(func() {
		interval, _ := fe.fields.Interval.Get()
		fe.model.IsDirty.Set(uint32(interval) != fe.frame.Interval)
	}))

	frm := container.New(layout.NewFormLayout(), ratelabel, fe.rateBox.Container)
	fe.Container = container.NewBorder(nil, fe.layerSelect, nil, nil, frm)
	fe.model.Frame.AddListener(binding.NewDataListener(fe.setFields))

	fe.tools = NewFrameTools(model, window)
	// sharedTools.AddItems(widget.NewToolbarSeparator())
	sharedTools.AddItems(fe.tools.Items()...)
	sharedTools.AddApply(fe.apply)
	sharedTools.AddRevert(fe.revert)

	return fe
}

func (fe *FrameEditor) setFields() {
	fe.frame = fe.model.GetFrame()
	fe.fields.FromFrame(fe.frame)
	fe.rateBox.Entry.SetText(strconv.FormatInt(int64(fe.frame.Interval), 10))
}

func (fe *FrameEditor) apply() {
	dirty, _ := fe.model.IsDirty.Get()
	if dirty {
		fe.fields.ToFrame(fe.frame)
		fe.model.UpdateFrame()
		fe.setFields()
	}
}

func (fe *FrameEditor) revert() {
	fe.setFields()
}
