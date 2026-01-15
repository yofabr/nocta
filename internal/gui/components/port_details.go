package components

import (
	"strings"

	"nocta/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PortDetails struct {
	container *fyne.Container
}

func NewPortDetails() *PortDetails {
	detailsTitle := widget.NewLabel("Port Details")
	detailsTitle.TextStyle = fyne.TextStyle{Bold: true}
	detailsTitle.Alignment = fyne.TextAlignCenter

	emptyStateLabel := widget.NewLabel("Select a port to see details")
	emptyStateLabel.Alignment = fyne.TextAlignCenter
	emptyStateLabel.Importance = widget.LowImportance

	detailsContent := container.NewVBox(
		detailsTitle,
		emptyStateLabel,
	)

	return &PortDetails{
		container: detailsContent,
	}
}

func (pd *PortDetails) UpdateDetails(port models.ActivePort, terminateCallback func()) {
	detailsTitle := widget.NewLabel("Port Details")
	detailsTitle.TextStyle = fyne.TextStyle{Bold: true}
	detailsTitle.Alignment = fyne.TextAlignCenter

	protocolLabel := widget.NewLabel("Protocol:")
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

	valueLabels := []*widget.Label{
		protocolValue, stateValue, addressValue, portValue, processValue,
		recvQValue, sendQValue, peerValue, pidValue,
		userValue, ppidValue, statValue, startedValue, elapsedValue,
	}
	for _, lbl := range valueLabels {
		lbl.Wrapping = fyne.TextWrapWord
		lbl.Importance = widget.MediumImportance
	}

	protocolValue.SetText(strings.ToUpper(port.Protocol))
	stateValue.SetText(port.State)
	addressValue.SetText(port.Addr)
	portValue.SetText(port.Port)
	processValue.SetText(port.Process)
	recvQValue.SetText(port.RecvQ)
	sendQValue.SetText(port.SendQ)
	peerValue.SetText(port.Peer_Addr_Port)
	pidValue.SetText(port.PID)
	userValue.SetText(port.PortDetails.User)
	ppidValue.SetText(port.PortDetails.PPID)
	statValue.SetText(port.PortDetails.STAT)
	startedValue.SetText(port.PortDetails.STARTED)
	elapsedValue.SetText(port.PortDetails.ELAPSED)

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
	)

	scrollableForm := container.NewScroll(container.NewPadded(detailsForm))
	scrollableForm.SetMinSize(fyne.NewSize(0, 400))

	objects := []fyne.CanvasObject{
		detailsTitle,
		widget.NewSeparator(),
		scrollableForm,
		widget.NewSeparator(),
	}

	if port.PID != "" {
		terminateBtn := widget.NewButton("Terminate", terminateCallback)
		terminateBtn.Importance = widget.DangerImportance
		actionButtons := container.NewGridWithColumns(3, terminateBtn)
		objects = append(objects, container.NewPadded(actionButtons))
	}

	pd.container.Objects = objects
	pd.container.Refresh()
}

func (pd *PortDetails) ClearDetails() {
	detailsTitle := widget.NewLabel("Port Details")
	detailsTitle.TextStyle = fyne.TextStyle{Bold: true}
	detailsTitle.Alignment = fyne.TextAlignCenter

	emptyStateLabel := widget.NewLabel("Select a port to see details")
	emptyStateLabel.Alignment = fyne.TextAlignCenter
	emptyStateLabel.Importance = widget.LowImportance

	pd.container.Objects = []fyne.CanvasObject{
		detailsTitle,
		emptyStateLabel,
	}
	pd.container.Refresh()
}

func (pd *PortDetails) GetContainer() *fyne.Container {
	return pd.container
}
