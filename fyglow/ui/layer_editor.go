package ui

import (
	"gglow/fyglow/effectio"
	"gglow/glow"
	"gglow/text"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type LayerEditor struct {
	*fyne.Container
	effect *effectio.EffectIo
	layer  *glow.Layer
	fields *effectio.LayerFields
	window fyne.Window

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

	imageLabel  *widget.Label
	imageButton *widget.Button

	isEditing bool
}

func NewLayerEditor(effect *effectio.EffectIo, window fyne.Window, menu *fyne.Menu) *LayerEditor {

	le := &LayerEditor{
		window:            window,
		effect:            effect,
		layer:             effect.GetCurrentLayer(),
		fields:            effectio.NewLayerFields(),
		rateBounds:        RateBounds,
		hueBounds:         HueShiftBounds,
		scanBounds:        ScanBounds,
		selectOrigin:      widget.NewSelect(text.OriginLabels, func(s string) {}),
		selectOrientation: widget.NewSelect(text.OrientationLabels, func(s string) {}),
	}

	le.createPatches()

	form := le.createForm()
	scroll := container.NewVScroll(form)
	tools := container.NewCenter(NewLayerToolbar())
	layerSelect := NewLayerSelect(effect)
	fixed := container.NewVBox(tools, layerSelect)
	le.Container = container.NewBorder(fixed, nil, nil, nil, scroll)

	listener := binding.NewDataListener(le.setFields)
	effect.AddFrameListener(listener)
	effect.AddLayerListener(listener)
	effect.OnSave(le.apply)
	return le
}

func (le *LayerEditor) setChanged() {
	if le.isEditing {
		le.effect.SetChanged()
	}
}

func (le *LayerEditor) createPatches() {
	le.patches = make([]*ColorPatch, effectio.MaxLayerColors)
	for i := 0; i < effectio.MaxLayerColors; i++ {
		patch := NewColorPatch()
		patch.SetOnTapped(le.selectColor(patch))
		patch.SetOnChanged(le.setChanged)
		le.patches[i] = patch
	}
}

func (le *LayerEditor) selectColor(patch *ColorPatch) func() {
	return func() {
		ce := NewColorPatchEditor(patch, le.window, le.setChanged)
		ce.Show()
	}
}

func (le *LayerEditor) createForm() *fyne.Container {
	labelOrigin := widget.NewLabel(text.OriginLabel.String())
	le.selectOrigin.OnChanged = func(s string) {
		current := le.layer.Grid.Origin
		selected := le.selectOrigin.SelectedIndex()
		if glow.Origin(selected) != current {
			le.fields.Origin.Set(selected)
			le.setChanged()
		}
	}

	labelOrientation := widget.NewLabel(text.OrientationLabel.String())
	le.selectOrientation.OnChanged = func(s string) {
		current := le.layer.Grid.Orientation
		selected := le.selectOrientation.SelectedIndex()
		if glow.Orientation(selected) != current {
			le.fields.Orientation.Set(selected)
			le.setChanged()
		}
	}

	scanLabel := widget.NewLabel(text.LengthLabel.String())
	scanCheckLabel := widget.NewLabel(text.ScanLabel.String())
	le.scanBox = NewRangeIntBox(le.fields.Scan, le.scanBounds)
	le.fields.Scan.AddListener(binding.NewDataListener(func() {
		scan, _ := le.fields.Scan.Get()
		if uint16(scan) != le.layer.Scan {
			le.setChanged()
		}
	}))
	le.checkScan = widget.NewCheck("", checkRangeBox(le.scanBox, le.fields.Scan))

	colorsLabel := widget.NewLabel(text.ColorsLabel.String())
	// gradientLabel := widget.NewLabel(resources.GradientLabel.String())

	huelabel := widget.NewLabel(text.HueShiftLabel.String())
	hueCheckLabel := widget.NewLabel(text.DynamicLabel.String())
	le.hueBox = NewRangeIntBox(le.fields.HueShift, le.hueBounds)
	le.fields.HueShift.AddListener(binding.NewDataListener(func() {
		shift, _ := le.fields.HueShift.Get()
		if int16(shift) != le.layer.HueShift {
			le.setChanged()
		}
	}))
	le.checkHue = widget.NewCheck("", checkRangeBox(le.hueBox, le.fields.HueShift))

	ratelabel := widget.NewLabel(text.RateLabel.String())
	rateCheckLabel := widget.NewLabel(text.OverrideLabel.String())
	le.rateBox = NewRangeIntBox(le.fields.Rate, le.rateBounds)
	le.fields.Rate.AddListener(binding.NewDataListener(func() {
		rate, _ := le.fields.Rate.Get()
		if uint32(rate) != le.layer.Rate {
			le.setChanged()
		}
	}))
	le.checkRate = widget.NewCheck("", checkRangeBox(le.rateBox, le.fields.Rate))

	patchBox := container.NewHBox()
	for _, patch := range le.patches {
		patchBox.Add(patch)
	}

	imageLoad := NewImageLoader(le.effect, le.window)
	le.imageLabel = widget.NewLabel(imageName(le.layer.ImageName))
	le.imageButton = widget.NewButton("Image...", func() {
		imageLoad.Start()
	})

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
		ratelabel, le.rateBox.Container,
		le.imageButton, le.imageLabel,
	)
	return frm
}

func (le *LayerEditor) setFields() {
	le.isEditing = false
	le.layer = le.effect.GetCurrentLayer()
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

	le.imageLabel.SetText(imageName(le.layer.ImageName))

	for i, p := range le.patches {
		if i < len(le.fields.Colors) {
			p.SetHSVColor(le.fields.Colors[i])
		} else {
			p.SetUnused(true)
		}
	}
	le.isEditing = true
}

func (le *LayerEditor) apply(frame *glow.Frame) {
	index := le.effect.LayerIndex()
	le.layer = frame.Layers[index]
	le.setColors()
	le.fields.ToLayer(le.layer)
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

func imageName(path string) (name string) {
	_, name = filepath.Split(path)
	ext := filepath.Ext(name)
	if len(ext) > 0 {
		name, _ = strings.CutSuffix(name, ext)
	}
	return
}
