package components

import (
	"time"

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
	lastRefresh       time.Time
	cooldown          time.Duration
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
		cooldown:       2 * time.Second,
	}
}

func (c *Controls) SetCallbacks(onSearchChanged func(string), onProtocolChanged func(string), onRefresh func()) {
	c.onSearchChanged = onSearchChanged
	c.onProtocolChanged = onProtocolChanged
	c.onRefresh = onRefresh

	c.searchEntry.OnChanged = onSearchChanged
	c.protocolSelect.OnChanged = onProtocolChanged

	throttledRefresh := func() {
		now := time.Now()
		if now.Sub(c.lastRefresh) < c.cooldown {
			return
		}
		c.lastRefresh = now
		onRefresh()
	}

	c.refreshButton.OnTapped = throttledRefresh
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
