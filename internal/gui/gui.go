package gui

import (
	"strings"

	"nocta/internal/application"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewGUI(appLogic application.Application) {
	a := app.New()
	w := a.NewWindow("Nocta")

	// -----------------------------
	// Dummy data
	// -----------------------------
	ports := appLogic.ActivePorts

	filtered := ports

	// -----------------------------
	// Right panel (details)
	// -----------------------------
	details := widget.NewLabel("Select a port to see details")
	details.Wrapping = fyne.TextWrapWord

	rightPanel := container.NewPadded(details)

	// -----------------------------
	// Ports list
	// -----------------------------
	list := widget.NewList(
		func() int {
			return len(filtered)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			p := filtered[i]
			o.(*widget.Label).SetText(
				strings.ToUpper(p.Protocol) + "  " + p.Addr + ":" + p.Port,
			)
		},
	)

	list.OnSelected = func(id widget.ListItemID) {
		p := filtered[id]

		details.SetText(
			"Protocol: " + p.Protocol + "\n" +
				"State: " + p.State + "\n" +
				"Address: " + p.Addr + "\n" +
				"Port: " + p.Port + "\n" +
				"Recv-Q: " + p.RecvQ + "\n" +
				"Send-Q: " + p.SendQ + "\n" +
				"Peer: " + p.Peer_Addr_Port + "\n" +
				"Process: " + p.Process,
		)
	}

	// -----------------------------
	// Search
	// -----------------------------
	search := widget.NewEntry()
	search.SetPlaceHolder("Search port, process...")

	search.OnChanged = func(s string) {
		filtered = nil
		for _, p := range ports {
			if strings.Contains(p.Port, s) ||
				strings.Contains(strings.ToLower(p.Process), strings.ToLower(s)) {
				filtered = append(filtered, p)
			}
		}
		list.UnselectAll()
		list.Refresh()
		details.SetText("Select a port to see details")
	}

	// -----------------------------
	// Protocol filter
	// -----------------------------
	protocolFilter := widget.NewSelect(
		[]string{"ALL", "TCP", "UDP"},
		func(value string) {
			filtered = nil
			for _, p := range ports {
				if value == "ALL" || strings.EqualFold(p.Protocol, value) {
					filtered = append(filtered, p)
				}
			}
			list.UnselectAll()
			list.Refresh()
			details.SetText("Select a port to see details")
		},
	)
	protocolFilter.SetSelected("ALL")

	// -----------------------------
	// Refresh button (dummy)
	// -----------------------------
	refreshBtn := widget.NewButton("⟳ Refresh", func() {
		// later: appLogic.RefreshPorts()
		list.Refresh()
	})

	// -----------------------------
	// Left panel layout
	// -----------------------------
	controls := container.NewVBox(
		search,
		container.NewGridWithColumns(2, protocolFilter, refreshBtn),
	)

	leftPanel := container.NewBorder(
		controls,
		nil,
		nil,
		nil,
		list,
	)

	// -----------------------------
	// Split layout
	// -----------------------------
	split := container.NewHSplit(leftPanel, rightPanel)
	split.Offset = 0.35

	w.SetContent(split)
	w.Resize(fyne.NewSize(900, 520))
	w.ShowAndRun()
}
