package gui

import (
	"fmt"
	"sync"

	"nocta/internal/config"
	"nocta/internal/gui/components"
	"nocta/internal/models"
	"nocta/internal/service"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

type App struct {
	service     service.PortService
	config      *config.Config
	portList    *components.PortList
	portDetails *components.PortDetails
	controls    *components.Controls
	window      *fyne.Window
	mu          sync.Mutex
	terminating bool
}

func NewGUI(portService service.PortService, cfg *config.Config) {
	a := &App{
		service: portService,
		config:  cfg,
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
	a.portDetails = components.NewPortDetails(*a.window)
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
	w.Resize(fyne.NewSize(float32(a.config.GUI.Width), float32(a.config.GUI.Height)))
}

func (a *App) loadInitialData() {
	ports, err := a.service.GetAllPorts()
	if err != nil {
		dialog.ShowError(fmt.Errorf("failed to load ports: %w", err), *a.window)
		return
	}
	a.portList.SetPorts(ports)
}

func (a *App) run() {
	(*a.window).ShowAndRun()
}

func (a *App) onPortSelected(port models.ActivePort) {
	if err := a.service.GetPortDetails(&port); err != nil {
		dialog.ShowError(fmt.Errorf("failed to get port details: %w", err), *a.window)
	}
	a.portDetails.UpdateDetails(port, func() {
		a.mu.Lock()
		if a.terminating {
			a.mu.Unlock()
			return
		}
		a.terminating = true
		a.mu.Unlock()

		defer func() {
			a.mu.Lock()
			a.terminating = false
			a.mu.Unlock()
		}()

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
		dialog.ShowError(fmt.Errorf("failed to refresh ports: %w", err), *a.window)
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
