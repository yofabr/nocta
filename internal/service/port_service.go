package service

import (
	"fmt"
	"nocta/internal/models"
	"nocta/internal/parser"
	"nocta/internal/scanner"
	"os/exec"
	"strconv"
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
	pid, err := validatePID(port.PID)
	if err != nil {
		return err
	}

	cmd := exec.Command("kill", "-TERM", strconv.Itoa(pid))
	return cmd.Run()
}

func (s *DefaultPortService) GetPortDetails(port *models.ActivePort) error {
	pid, err := validatePID(port.PID)
	if err != nil {
		return err
	}

	cmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "user,pid,ppid,stat,etime,lstart,comm,args")
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
	parser.SetLabels(&labels, label)

	portDetails := []models.PortDetail{}
	setPortDetails(&portDetails, rows[1:], labels)

	if len(portDetails) > 0 {
		port.PortDetails = portDetails[0]
	}

	return nil
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

		*portDetails = append(*portDetails, parser.MapPortDetails(fields, labels))
	}
}

func validatePID(pid string) (int, error) {
	trimmed := strings.TrimSpace(pid)
	if trimmed == "" {
		return 0, fmt.Errorf("missing pid")
	}

	parsedPID, err := strconv.Atoi(trimmed)
	if err != nil || parsedPID <= 0 {
		return 0, fmt.Errorf("invalid pid: %q", pid)
	}

	return parsedPID, nil
}
