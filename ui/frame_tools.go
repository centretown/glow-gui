package ui

import (
	"gglow/iohandler"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type FrameTools struct {
	frameMenu *ButtonItem

	newFrame     *ButtonItem
	newFolder    *ButtonItem
	deleteFrame  *ButtonItem
	createDialog *EffectDialog
	folderDialog *FolderDialog

	toolBar *widget.Toolbar
	popUp   *widget.PopUp
	effect  iohandler.EffectIoHandler
}

func NewFrameTools(effect iohandler.EffectIoHandler, window fyne.Window) *FrameTools {
	ft := &FrameTools{
		effect: effect,
	}

	ft.createDialog = NewEffectDialog(effect, window)
	ft.folderDialog = NewFolderDialog(effect, window)

	ft.frameMenu = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentIcon(), ft.menu))

	ft.newFolder = NewButtonItem(
		widget.NewButtonWithIcon("", theme.FolderNewIcon(), func() {
			ft.popUp.Hide()
			ft.folderDialog.Start()
		}))
	ft.newFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DocumentCreateIcon(), func() {
			ft.popUp.Hide()
			ft.createDialog.Start()
		}))
	ft.deleteFrame = NewButtonItem(
		widget.NewButtonWithIcon("", theme.DeleteIcon(), ft.delete))
	ft.toolBar = widget.NewToolbar(
		ft.newFolder,
		ft.newFrame,
		ft.deleteFrame,
	)

	ft.popUp = widget.NewPopUp(ft.toolBar, window.Canvas())
	return ft
}

func (ft *FrameTools) Items() (items []widget.ToolbarItem) {
	items = []widget.ToolbarItem{
		ft.frameMenu,
	}
	return
}

func (ft *FrameTools) delete() {
	ft.popUp.Hide()
}

func (ft *FrameTools) menu() {
	pos := fyne.Position{
		X: -ft.popUp.MinSize().Width/2 + ft.frameMenu.Button.MinSize().Width/2,
		Y: -ft.popUp.MinSize().Height,
	}

	ft.popUp.ShowAtRelativePosition(pos, ft.frameMenu.Button)
}
