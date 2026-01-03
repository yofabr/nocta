package gui

import (
	"strings"

	"nocta/internal/application"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewGUI(appLogic *application.Application) {
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
	detailsTitle := widget.NewLabel("Port Details")
	detailsTitle.TextStyle = fyne.TextStyle{Bold: true}
	detailsTitle.Alignment = fyne.TextAlignCenter

	// Detail fields
	protocolLabel := widget.NewLabel("Protocol (Net-id):")
	protocolValue := widget.NewLabel("")
	stateLabel := widget.NewLabel("State:")
	stateValue := widget.NewLabel("")
	addressLabel := widget.NewLabel("Address:")
	addressValue := widget.NewLabel("")
	portLabel := widget.NewLabel("Port:")
	portValue := widget.NewLabel("")
	processLabel := widget.NewLabel("Process:")
	processValue := widget.NewLabel("")
	recvQLabel := widget.NewLabel("Recv-Q:")
	recvQValue := widget.NewLabel("")
	sendQLabel := widget.NewLabel("Send-Q:")
	sendQValue := widget.NewLabel("")
	peerLabel := widget.NewLabel("Peer:")
	peerValue := widget.NewLabel("")
	pidLabel := widget.NewLabel("PID:")
	pidValue := widget.NewLabel("")
	userLabel := widget.NewLabel("USER:")
	userValue := widget.NewLabel("")
	ppidLabel := widget.NewLabel("PPID:")
	ppidValue := widget.NewLabel("")
	statLabel := widget.NewLabel("STAT:")
	statValue := widget.NewLabel("")
	startedLabel := widget.NewLabel("STARTED:")
	startedValue := widget.NewLabel("")
	elapsedLabel := widget.NewLabel("ELAPSED:")
	elapsedValue := widget.NewLabel("")
	commandLabel := widget.NewLabel("COMMAND:")
	commandValue := widget.NewLabel("")
	// Make all value labels wrap text and set importance
	valueLabels := []*widget.Label{
		protocolValue, stateValue, addressValue, portValue, processValue,
		recvQValue, sendQValue, peerValue, pidValue,
		userValue, ppidValue, statValue, startedValue, elapsedValue, commandValue,
	}
	for _, lbl := range valueLabels {
		lbl.Wrapping = fyne.TextWrapWord
		lbl.Importance = widget.MediumImportance
	}

	// Details form
	detailsForm := container.NewGridWithColumns(2,
		protocolLabel, protocolValue,
		stateLabel, stateValue,
		addressLabel, addressValue,
		portLabel, portValue,
		processLabel, processValue,
		recvQLabel, recvQValue,
		sendQLabel, sendQValue,
		peerLabel, peerValue,
		pidLabel, pidValue,
		userLabel, userValue,
		ppidLabel, ppidValue,
		statLabel, statValue,
		startedLabel, startedValue,
		elapsedLabel, elapsedValue,
		commandLabel, commandValue,
	)

	killBtn := widget.NewButton("Kill", func() {
		// TODO: Implement delete functionality
	})
	killBtn.Importance = widget.DangerImportance

	actionButtons := container.NewGridWithColumns(3, killBtn)

	// Empty state message
	emptyStateLabel := widget.NewLabel("Select a port to see details")
	emptyStateLabel.Alignment = fyne.TextAlignCenter
	emptyStateLabel.Importance = widget.LowImportance

	// Details container - will switch between empty state and details
	detailsContent := container.NewVBox(
		detailsTitle,
		emptyStateLabel,
	)

	// Update details function
	updateDetails := func(p application.ActivePort) {
		p.Detail()
		protocolValue.SetText(strings.ToUpper(p.Protocol))
		stateValue.SetText(p.State)
		addressValue.SetText(p.Addr)
		portValue.SetText(p.Port)
		processValue.SetText(p.Process)
		recvQValue.SetText(p.RecvQ)
		sendQValue.SetText(p.SendQ)
		peerValue.SetText(p.Peer_Addr_Port)
		pidValue.SetText(p.PID)
		userValue.SetText(p.PortDetails.User)
		ppidValue.SetText(p.PortDetails.PPID)
		statValue.SetText(p.PortDetails.STAT)
		startedValue.SetText(p.PortDetails.STARTED)
		elapsedValue.SetText(p.PortDetails.ELAPSED)
		commandValue.SetText(p.PortDetails.COMMAND)

		// Create scrollable container for the details form
		scrollableForm := container.NewScroll(container.NewPadded(detailsForm))
		scrollableForm.SetMinSize(fyne.NewSize(0, 400))

		detailsContent.Objects = []fyne.CanvasObject{
			detailsTitle,
			widget.NewSeparator(),
			scrollableForm,
			widget.NewSeparator(),
			container.NewPadded(actionButtons),
		}
		detailsContent.Refresh()
	}

	clearDetails := func() {
		detailsContent.Objects = []fyne.CanvasObject{
			detailsTitle,
			emptyStateLabel,
		}
		detailsContent.Refresh()
	}

	rightPanel := container.NewPadded(
		container.NewBorder(
			nil,
			nil,
			nil,
			nil,
			detailsContent,
		),
	)

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
		updateDetails(p)
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
		clearDetails()
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
			clearDetails()
		},
	)
	protocolFilter.SetSelected("ALL")

	// -----------------------------
	// Refresh button
	// -----------------------------
	refreshBtn := widget.NewButton("⟳ Refresh", func() {
		appLogic.RefreshPorts()
		ports = appLogic.ActivePorts

		// Reapply filters
		searchText := search.Text
		protocolValue := protocolFilter.Selected
		filtered = nil

		for _, p := range ports {
			// Apply protocol filter
			if protocolValue != "ALL" && !strings.EqualFold(p.Protocol, protocolValue) {
				continue
			}
			// Apply search filter
			if searchText != "" {
				if !strings.Contains(p.Port, searchText) &&
					!strings.Contains(strings.ToLower(p.Process), strings.ToLower(searchText)) {
					continue
				}
			}
			filtered = append(filtered, p)
		}

		list.UnselectAll()
		list.Refresh()
		clearDetails()
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
