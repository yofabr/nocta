package service

import (
	"nocta/internal/models"
	"nocta/internal/scanner"
	"os/exec"
	"strings"
)

type PortService interface {
	GetAllPorts() ([]models.ActivePort, error)
	RefreshPorts() ([]models.ActivePort, error)
	QueryPort(query models.QueryParams) ([]models.ActivePort, error)
	TerminatePort(port models.ActivePort) error
	GetPortDetails(port *models.ActivePort) error
}

type DefaultPortService struct {
	scanner scanner.Scanner
}

func NewPortService() *DefaultPortService {
	return &DefaultPortService{
		scanner: scanner.NewSystemScanner(),
	}
}

func (s *DefaultPortService) GetAllPorts() ([]models.ActivePort, error) {
	return s.scanner.ListAllPorts()
}

func (s *DefaultPortService) RefreshPorts() ([]models.ActivePort, error) {
	return s.GetAllPorts()
}

func (s *DefaultPortService) QueryPort(query models.QueryParams) ([]models.ActivePort, error) {
	return s.scanner.QueryPort(query)
}

func (s *DefaultPortService) TerminatePort(port models.ActivePort) error {
	if port.PID == "" {
		return nil
	}

	cmd := exec.Command("kill", port.PID)
	return cmd.Run()
}

func (s *DefaultPortService) GetPortDetails(port *models.ActivePort) error {
	cmd := exec.Command("ps", "-p", port.PID, "-o", "user,pid,ppid,stat,etime,lstart,comm,args")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	rows := strings.Split(string(output), "\n")
	if len(rows) < 2 {
		return nil
	}

	label := rows[0]
	labels := map[int]string{}
	setLabels(&labels, label)

	portDetails := []models.PortDetail{}
	setPortDetails(&portDetails, rows[1:], labels)

	if len(portDetails) > 0 {
		port.PortDetails = portDetails[0]
	}

	return nil
}

func setLabels(labels *map[int]string, label string) {
	fields := strings.Fields(strings.TrimSpace(label))
	for i, l := range fields {
		(*labels)[i] = l
	}
}

func setPortDetails(portDetails *[]models.PortDetail, lines []string, labels map[int]string) {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		*portDetails = append(*portDetails, mapPortDetails(fields, labels))
	}
}

func mapPortDetails(port []string, labels map[int]string) models.PortDetail {
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
