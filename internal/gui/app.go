package gui

import (
	"nocta/internal/gui/components"
	"nocta/internal/models"
	"nocta/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

type App struct {
	service     service.PortService
	portList    *components.PortList
	portDetails *components.PortDetails
	controls    *components.Controls
	window      *fyne.Window
}

func NewGUI(portService service.PortService) {
	a := &App{
		service: portService,
	}

	a.setupUI()
	a.loadInitialData()
	a.run()
}

func (a *App) setupUI() {
	fyneApp := app.New()
	w := fyneApp.NewWindow("Nocta")
	a.window = &w

	a.portList = components.NewPortList(a.onPortSelected)
	a.portDetails = components.NewPortDetails()
	a.controls = components.NewControls()

	a.controls.SetCallbacks(
		a.onSearchChanged,
		a.onProtocolChanged,
		a.onRefresh,
	)

	leftPanel := container.NewBorder(
		a.controls.GetContainer(),
		nil,
		nil,
		nil,
		a.portList.GetList(),
	)

	rightPanel := container.NewPadded(
		container.NewBorder(
			nil,
			nil,
			nil,
			nil,
			a.portDetails.GetContainer(),
		),
	)

	split := container.NewHSplit(leftPanel, rightPanel)
	split.Offset = 0.35

	w.SetContent(split)
	w.Resize(fyne.NewSize(900, 520))
}

func (a *App) loadInitialData() {
	ports, err := a.service.GetAllPorts()
	if err != nil {
		return
	}
	a.portList.SetPorts(ports)
}

func (a *App) run() {
	(*a.window).ShowAndRun()
}

func (a *App) onPortSelected(port models.ActivePort) {
	a.service.GetPortDetails(&port)
	a.portDetails.UpdateDetails(port, func() {
		a.service.TerminatePort(port)
		a.onRefresh()
	})
}

func (a *App) onSearchChanged(searchText string) {
	a.applyFilters()
}

func (a *App) onProtocolChanged(protocol string) {
	a.applyFilters()
}

func (a *App) onRefresh() {
	ports, err := a.service.RefreshPorts()
	if err != nil {
		return
	}
	a.portList.SetPorts(ports)
	a.applyFilters()
}

func (a *App) applyFilters() {
	searchText := a.controls.GetSearchText()
	protocol := a.controls.GetSelectedProtocol()
	a.portList.Filter(searchText, protocol)
}
