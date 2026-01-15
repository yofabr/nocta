package scanner

import (
	"nocta/internal/models"
	"os/exec"
	"regexp"
	"strings"
)

type Scanner interface {
	ListAllPorts() ([]models.ActivePort, error)
	QueryPort(query models.QueryParams) ([]models.ActivePort, error)
}

type SystemScanner struct{}

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
	setLabels(&labels, label)

	setActivePorts(&activePorts, cmdList[1:])
	return activePorts
}

func extractPID(s string) string {
	re := regexp.MustCompile(`pid=(\d+)`)
	matches := re.FindStringSubmatch(s)
	if len(matches) == 0 {
		return ""
	}
	return matches[1]
}

func setLabels(labels *map[int]string, label string) {
	fields := strings.Fields(strings.TrimSpace(label))
	for i, l := range fields {
		(*labels)[i] = l
	}
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
