package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Controls struct {
	searchEntry       *widget.Entry
	protocolSelect    *widget.Select
	refreshButton     *widget.Button
	container         *fyne.Container
	onSearchChanged   func(string)
	onProtocolChanged func(string)
	onRefresh         func()
}

func NewControls() *Controls {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search port, process...")

	protocolSelect := widget.NewSelect(
		[]string{"ALL", "TCP", "UDP"},
		nil,
	)
	protocolSelect.SetSelected("ALL")

	refreshButton := widget.NewButton("⟳ Refresh", nil)

	controls := container.NewVBox(
		searchEntry,
		container.NewGridWithColumns(2, protocolSelect, refreshButton),
	)

	return &Controls{
		searchEntry:    searchEntry,
		protocolSelect: protocolSelect,
		refreshButton:  refreshButton,
		container:      controls,
	}
}

func (c *Controls) SetCallbacks(onSearchChanged func(string), onProtocolChanged func(string), onRefresh func()) {
	c.onSearchChanged = onSearchChanged
	c.onProtocolChanged = onProtocolChanged
	c.onRefresh = onRefresh

	c.searchEntry.OnChanged = onSearchChanged
	c.protocolSelect.OnChanged = onProtocolChanged
	c.refreshButton.OnTapped = onRefresh
}

func (c *Controls) GetSearchText() string {
	return c.searchEntry.Text
}

func (c *Controls) GetSelectedProtocol() string {
	return c.protocolSelect.Selected
}

func (c *Controls) GetContainer() *fyne.Container {
	return c.container
}
