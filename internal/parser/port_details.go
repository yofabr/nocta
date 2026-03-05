package parser

import (
	"nocta/internal/models"
	"strings"
)

func SetLabels(labels *map[int]string, label string) {
	fields := strings.Fields(strings.TrimSpace(label))
	for i, l := range fields {
		(*labels)[i] = l
	}
}

func MapPortDetails(port []string, labels map[int]string) models.PortDetail {
	if len(port) == 0 {
		return models.PortDetail{}
	}

	userIdx, ppidIdx, statIdx, elapsedIdx, startedIdx := -1, -1, -1, -1, -1

	for idx, label := range labels {
		switch strings.ToUpper(label) {
		case "USER":
			userIdx = idx
		case "PPID":
			ppidIdx = idx
		case "STAT":
			statIdx = idx
		case "ELAPSED":
			elapsedIdx = idx
		case "STARTED":
			startedIdx = idx
		}
	}

	getField := func(idx int) string {
		if idx >= 0 && idx < len(port) {
			return port[idx]
		}
		return ""
	}

	started := ""
	if startedIdx >= 0 && startedIdx+4 < len(port) {
		startedParts := []string{
			port[startedIdx],
			port[startedIdx+1],
			port[startedIdx+2],
			port[startedIdx+3],
			port[startedIdx+4],
		}
		started = strings.Join(startedParts, " ")
	}

	return models.PortDetail{
		User:    getField(userIdx),
		PPID:    getField(ppidIdx),
		STAT:    getField(statIdx),
		ELAPSED: getField(elapsedIdx),
		STARTED: started,
	}
}
