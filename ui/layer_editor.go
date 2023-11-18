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

type LayerEditor struct {
	*fyne.Container
	model   *data.Model
	layer   *glow.Layer
	fields  *data.LayerFields
	window  fyne.Window
	isDirty binding.Bool

	patches []*ColorPatch

	bDynamic  bool
	bScan     bool
	bOverride bool

	selectOrigin      *widget.Select
	selectOrientation *widget.Select

	checkScan *widget.Check
	checkHue  *widget.Check
	checkRate *widget.Check

	scanBox *RangeIntBox
	hueBox  *RangeIntBox
	rateBox *RangeIntBox

	rateBounds *IntEntryBounds
	hueBounds  *IntEntryBounds
	scanBounds *IntEntryBounds

	tools *LayerTools
}

func NewLayerEditor(model *data.Model, isDirty binding.Bool, window fyne.Window,
	sharedTools *SharedTools) *LayerEditor {

	le := &LayerEditor{
		window: window,

		model: model,
		layer: model.GetCurrentLayer(),

		fields:  data.NewLayerFields(),
		isDirty: isDirty,
		tools:   NewLayerTools(model),

		rateBounds: RateBounds,
		hueBounds:  HueShiftBounds,
		scanBounds: ScanBounds,

		selectOrigin:      widget.NewSelect(resources.OriginLabels, func(s string) {}),
		selectOrientation: widget.NewSelect(resources.OrientationLabels, func(s string) {}),
	}

	le.createPatches()

	form := le.createForm()
	scroll := container.NewVScroll(form)

	le.Container = container.NewBorder(nil, nil, nil, nil, scroll)

	le.model.Layer.AddListener(binding.NewDataListener(le.setFields))

	sharedTools.AddItems(le.tools.Items()...)
	sharedTools.AddApply(le.apply)
	sharedTools.AddRevert(le.revert)
	return le
}

func (le *LayerEditor) createPatches() {
	le.patches = make([]*ColorPatch, data.MaxLayerColors)
	for i := 0; i < data.MaxLayerColors; i++ {
		patch := NewColorPatch(le.isDirty)
		patch.SetTapped(le.selectColor(patch))
		le.patches[i] = patch
	}
}

func (le *LayerEditor) selectColor(patch *ColorPatch) func() {
	return func() {
		ce := NewColorPatchEditor(patch, le.isDirty, le.window)
		ce.Show()
	}
}

func (le *LayerEditor) createForm() *fyne.Container {
	labelOrigin := widget.NewLabel(resources.OriginLabel.String())
	le.selectOrigin.OnChanged = func(s string) {
		current := le.layer.Grid.Origin
		selected := le.selectOrigin.SelectedIndex()
		if glow.Origin(selected) != current {
			le.fields.Origin.Set(selected)
			le.isDirty.Set(true)
		}
	}

	labelOrientation := widget.NewLabel(resources.OrientationLabel.String())
	le.selectOrientation.OnChanged = func(s string) {
		current := le.layer.Grid.Orientation
		selected := le.selectOrientation.SelectedIndex()
		if glow.Orientation(selected) != current {
			le.fields.Orientation.Set(selected)
			le.isDirty.Set(true)
		}
	}

	scanLabel := widget.NewLabel(resources.LengthLabel.String())
	scanCheckLabel := widget.NewLabel(resources.ScanLabel.String())
	le.scanBox = NewRangeIntBox(le.fields.Scan, le.scanBounds)
	le.fields.Scan.AddListener(binding.NewDataListener(func() {
		scan, _ := le.fields.Scan.Get()
		le.isDirty.Set(uint16(scan) != le.layer.Scan)
	}))
	le.checkScan = widget.NewCheck("", checkRangeBox(le.scanBox, le.fields.Scan))

	colorsLabel := widget.NewLabel(resources.ColorsLabel.String())
	// gradientLabel := widget.NewLabel(resources.GradientLabel.String())

	huelabel := widget.NewLabel(resources.HueShiftLabel.String())
	hueCheckLabel := widget.NewLabel(resources.DynamicLabel.String())
	le.hueBox = NewRangeIntBox(le.fields.HueShift, le.hueBounds)
	le.fields.HueShift.AddListener(binding.NewDataListener(func() {
		shift, _ := le.fields.HueShift.Get()
		le.isDirty.Set(int16(shift) != le.layer.HueShift)
	}))
	le.checkHue = widget.NewCheck("", checkRangeBox(le.hueBox, le.fields.HueShift))

	ratelabel := widget.NewLabel(resources.RateLabel.String())
	rateCheckLabel := widget.NewLabel(resources.OverrideLabel.String())
	le.rateBox = NewRangeIntBox(le.fields.Rate, le.rateBounds)
	le.fields.Rate.AddListener(binding.NewDataListener(func() {
		rate, _ := le.fields.Rate.Get()
		le.isDirty.Set(uint32(rate) != le.layer.Rate)
	}))
	le.checkRate = widget.NewCheck("", checkRangeBox(le.rateBox, le.fields.Rate))

	patchBox := container.NewHBox()
	for _, patch := range le.patches {
		patchBox.Add(patch)
	}

	sep := widget.NewSeparator()
	frm := container.New(layout.NewFormLayout(),
		sep, sep,
		labelOrigin, le.selectOrigin,
		labelOrientation, le.selectOrientation,
		scanCheckLabel, le.checkScan,
		scanLabel, le.scanBox.Container,
		sep, sep,
		colorsLabel, patchBox,
		hueCheckLabel, le.checkHue,
		huelabel, le.hueBox.Container,
		sep, sep,
		rateCheckLabel, le.checkRate,
		ratelabel, le.rateBox.Container)
	return frm
}

func (le *LayerEditor) setFields() {
	le.layer = le.model.GetCurrentLayer()
	le.fields.FromLayer(le.layer)

	le.selectOrigin.SetSelectedIndex(int(le.layer.Grid.Origin))
	le.selectOrientation.SetSelectedIndex(int(le.layer.Grid.Orientation))

	le.bDynamic = (le.layer.HueShift != int16(le.hueBounds.OffVal))
	le.hueBox.Entry.SetText(strconv.FormatInt(int64(le.layer.HueShift), 10))
	le.checkHue.SetChecked(le.bDynamic)
	le.hueBox.Enable(le.bDynamic)

	le.bScan = (le.layer.Scan != uint16(le.scanBounds.OffVal))
	le.scanBox.Entry.SetText(strconv.FormatInt(int64(le.layer.Scan), 10))
	le.checkScan.SetChecked(le.bScan)
	le.scanBox.Enable(le.bScan)

	le.bOverride = (le.layer.Rate != uint32(le.rateBounds.OffVal))
	le.checkRate.SetChecked(le.bOverride)
	le.rateBox.Entry.SetText(strconv.FormatInt(int64(le.layer.Rate), 10))
	le.rateBox.Enable(le.bOverride)

	for i, p := range le.patches {
		if i < len(le.fields.Colors) {
			p.SetHSVColor(le.fields.Colors[i])
		} else {
			p.SetUnused(true)
		}
	}
}

func (le *LayerEditor) apply() {
	dirty, _ := le.isDirty.Get()
	if dirty {
		le.setColors()
		le.fields.ToLayer(le.layer)
		le.model.UpdateFrame()
		le.setFields()
	}
}

func (le *LayerEditor) revert() {
	le.setFields()
}

func (le *LayerEditor) setColors() {
	var colors []glow.HSV = make([]glow.HSV, 0)
	for _, p := range le.patches {
		if !p.Unused() {
			colors = append(colors, p.GetHSVColor())
		}
	}
	le.fields.Colors = colors
}
