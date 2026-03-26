package scanner

import (
	"errors"
	"nocta/internal/models"
	"nocta/internal/parser"
	"os/exec"
	"regexp"
	"strings"
)

type Scanner interface {
	ListAllPorts() ([]models.ActivePort, error)
	QueryPort(query models.QueryParams) ([]models.ActivePort, error)
}

type SystemScanner struct{}

var pidRegex = regexp.MustCompile(`pid=(\d+)`)

func NewSystemScanner() *SystemScanner {
	return &SystemScanner{}
}

func (s *SystemScanner) ListAllPorts() ([]models.ActivePort, error) {
	cmd := exec.Command("ss", "-tulnp")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parsePortString(string(out)), nil
}

func (s *SystemScanner) QueryPort(query models.QueryParams) ([]models.ActivePort, error) {
	protocol, err := query.ValidateProtocol()
	if err != nil {
		return nil, err
	}

	port := query.ValidatePort()
	if port == "" {
		return nil, errors.New("invalid port number (must be 1-65535)")
	}
	filter := strings.Builder{}
	filter.WriteString("sport = :")
	filter.WriteString(port)

	cmd := exec.Command(
		"ss",
		protocol,
		"-l",
		"-n",
		"-p",
		filter.String(),
	)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parsePortString(string(out)), nil
}

func parsePortString(portString string) []models.ActivePort {
	var activePorts []models.ActivePort

	cmdList := strings.Split(portString, "\n")
	if len(cmdList) < 2 {
		return activePorts
	}

	label := cmdList[0]
	labels := map[int]string{}
	parser.SetLabels(&labels, label)

	setActivePorts(&activePorts, cmdList[1:])
	return activePorts
}

func extractPID(s string) string {
	matches := pidRegex.FindStringSubmatch(s)
	if len(matches) == 0 {
		return ""
	}
	return matches[1]
}

func parseAddrPort(s string) (addr, port string) {
	lastColon := strings.LastIndex(s, ":")
	if lastColon == -1 {
		return s, ""
	}

	addr = s[:lastColon]
	port = s[lastColon+1:]
	return
}

func mapPort(port []string) models.ActivePort {
	if len(port) < 7 {
		return models.ActivePort{}
	}

	addr, prt := parseAddrPort(port[4])
	activePort := models.ActivePort{
		Protocol:       port[0],
		State:          port[1],
		RecvQ:          port[2],
		SendQ:          port[3],
		Addr:           addr,
		Port:           prt,
		Process:        port[6],
		PID:            extractPID(port[6]),
		Peer_Addr_Port: port[5],
	}

	return activePort
}

func setActivePorts(activePorts *[]models.ActivePort, lines []string) {
	for _, line := range lines {
		fields := strings.Fields(line)

		if len(fields) > 7 {
			merged := make([]string, 7)
			copy(merged, fields[:6])
			merged[6] = strings.Join(fields[6:], " ")
			fields = merged
		} else if len(fields) < 7 {
			for len(fields) < 7 {
				fields = append(fields, "")
			}
		}

		*activePorts = append(*activePorts, mapPort(fields))
	}
}
