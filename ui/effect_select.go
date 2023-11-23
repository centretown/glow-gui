package ui

import (
	"fmt"
	"glow-gui/control"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type EffectSelect struct {
	*widget.Select
	model *control.Manager
	auto  bool
}

func NewEffectSelect(model *control.Manager) *widget.Select {
	fs := &EffectSelect{
		model: model,
	}
	fs.Select = widget.NewSelect(model.KeyList(), fs.onChange)
	model.AddFrameListener(binding.NewDataListener(func() {
		selected := fs.Select.Selected
		if selected != model.EffectName() {
			fs.auto = true
			fs.Select.SetSelected(model.EffectName())
		}
	}))
	return fs.Select
}

func (fs *EffectSelect) onChange(title string) {
	if fs.model.IsFolder(title) {
		fs.auto = false
		fs.Select.SetOptions(fs.model.RefreshKeys(title))
		return
	}
	if fs.auto {
		fs.auto = false
	} else {
		fs.model.ReadEffect(title)
		fmt.Println("EffectSelect", title)
	}
}
