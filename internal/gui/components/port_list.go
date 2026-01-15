package components

import (
	"strings"

	"nocta/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type PortList struct {
	list       *widget.List
	ports      []models.ActivePort
	filtered   []models.ActivePort
	onSelected func(models.ActivePort)
}

func NewPortList(onSelected func(models.ActivePort)) *PortList {
	pl := &PortList{
		onSelected: onSelected,
	}

	pl.list = widget.NewList(
		func() int {
			return len(pl.filtered)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			port := pl.filtered[i]
			o.(*widget.Label).SetText(
				strings.ToUpper(port.Protocol) + "  " + port.Addr + ":" + port.Port,
			)
		},
	)

	pl.list.OnSelected = func(id widget.ListItemID) {
		if id >= 0 && id < len(pl.filtered) {
			pl.onSelected(pl.filtered[id])
		}
	}

	return pl
}

func (pl *PortList) SetPorts(ports []models.ActivePort) {
	pl.ports = ports
	pl.filtered = ports
	pl.list.Refresh()
}

func (pl *PortList) Filter(searchText, protocolFilter string) {
	pl.filtered = nil

	for _, port := range pl.ports {
		protocolMatch := protocolFilter == "ALL" || strings.EqualFold(port.Protocol, protocolFilter)

		var searchMatch bool
		if searchText == "" {
			searchMatch = true
		} else {
			searchMatch = strings.Contains(port.Port, searchText) ||
				strings.Contains(strings.ToLower(port.Process), strings.ToLower(searchText))
		}

		if protocolMatch && searchMatch {
			pl.filtered = append(pl.filtered, port)
		}
	}

	pl.list.UnselectAll()
	pl.list.Refresh()
}

func (pl *PortList) GetList() *widget.List {
	return pl.list
}

func (pl *PortList) GetFilteredPorts() []models.ActivePort {
	return pl.filtered
}
