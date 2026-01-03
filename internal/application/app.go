package application

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type PortDetail struct {
	User string
	// PID     string
	// PID moved out of here...
	PPID    string
	STAT    string
	ELAPSED string
	STARTED string
}

type ActivePort struct {
	Protocol       string
	State          string
	Addr           string
	Port           string
	RecvQ          string
	SendQ          string
	Peer_Addr_Port string
	Process        string
	PortDetails    PortDetail
	PID            string
}

type Application struct {
	ActivePorts []ActivePort
}

func NewApplication() *Application {
	app := Application{}
	port_string := app.ListAllPorts()

	active_ports := formatPortString(port_string)
	app.setActivePorts(active_ports)
	return &app
}

func (app *Application) RefreshPorts() {
	port_string := app.ListAllPorts()

	active_ports := formatPortString(port_string)
	app.setActivePorts(active_ports)
}

func (a *Application) setActivePorts(active_ports []ActivePort) {
	a.ActivePorts = active_ports
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
	// Keep brackets, just find the last colon
	lastColon := strings.LastIndex(s, ":")
	if lastColon == -1 {
		// No port found
		return s, ""
	}

	addr = s[:lastColon]   // everything before last colon
	port = s[lastColon+1:] // everything after last colon
	return
}

func mapPort(port []string) ActivePort {
	if len(port) < 7 {
		return ActivePort{}
	}

	addr, prt := parseAddrPort(port[4])
	active_port := ActivePort{
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

	return active_port
}

func mapPortDetails(port []string, labels map[int]string) PortDetail {
	if len(port) == 0 {
		return PortDetail{}
	}

	// Find column indices from labels
	userIdx := -1
	ppidIdx := -1
	statIdx := -1
	elapsedIdx := -1
	startedIdx := -1

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

	// Helper to safely get field value
	getField := func(idx int) string {
		if idx >= 0 && idx < len(port) {
			return port[idx]
		}
		return ""
	}

	// Parse STARTED field - it spans multiple columns (date + time)
	// STARTED format: "Sat Jan  3 14:31:26 2026" (5 fields in data, 1 column in header)
	// The ps lstart format always produces 5 fields: day, month, day-of-month, time, year
	started := ""
	if startedIdx >= 0 && startedIdx+4 < len(port) {
		// STARTED always spans exactly 5 fields in the data row
		startedParts := []string{
			port[startedIdx],   // day (e.g., "Sat")
			port[startedIdx+1], // month (e.g., "Jan")
			port[startedIdx+2], // day-of-month (e.g., "3")
			port[startedIdx+3], // time (e.g., "14:31:26")
			port[startedIdx+4], // year (e.g., "2026")
		}
		started = strings.Join(startedParts, " ")
	}

	// COMMAND is the last field (full args - the second COMMAND column)
	// The last COMMAND column (args) is always the last field in the data row

	port_detail := PortDetail{
		User:    getField(userIdx),
		PPID:    getField(ppidIdx),
		STAT:    getField(statIdx),
		ELAPSED: getField(elapsedIdx),
		STARTED: started,
	}

	return port_detail
}

func setActivePorts(_activePorts *[]ActivePort, lines []string) {
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

		*(_activePorts) = append(*_activePorts, mapPort(fields))
	}
}

func setPortDetails(_portDetails *[]PortDetail, lines []string, labels map[int]string) {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		*(_portDetails) = append(*_portDetails, mapPortDetails(fields, labels))
	}
}

func formatPortString(port_string string) []ActivePort {
	var activePorts []ActivePort

	cmd_list := strings.Split(port_string, "\n")

	label := cmd_list[0]
	labels := map[int]string{}
	setLabels(&labels, label)

	setActivePorts(&activePorts, cmd_list[1:])

	return activePorts
}

func (a *Application) ListAllPorts() string {
	cmd := exec.Command("ss", "-tulnp")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

type QueryParams struct {
	Port     int
	Protocol string
}

func (q *QueryParams) ValidatePort() string {
	return strconv.Itoa(q.Port)
}

func (q *QueryParams) ValidateProtocol() (string, error) {
	pr := strings.ToLower(strings.TrimSpace(q.Protocol))

	switch pr {
	case "tcp":
		return "-t", nil
	case "udp":
		return "-u", nil
	case "both", "":
		return "-tu", nil
	default:
		return "", errors.New("invalid protocol (tcp, udp, both)")
	}
}

func (a *Application) QueryPort(query QueryParams) string {
	protocol, err := query.ValidateProtocol()
	if err != nil {
		log.Fatal(err)
	}

	port := strconv.Itoa(query.Port)
	filter := fmt.Sprintf("sport = :%s", port)

	cmd := exec.Command(
		"ss",
		protocol, // "-t", "-u", or "-tu"
		"-l",
		"-n",
		"-p",
		filter,
	)

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	return string(out)
}

func (a *ActivePort) KillPort() {
	// This is where you terminate the running port

	if a.PID == "" {
		return
	}

	cmd := exec.Command("kill", a.PID)
	cmd.Run()

	// fmt.Println("Kill message:", output)
}

func (a *ActivePort) Detail() {
	// a.PortDetails = PortDetail{
	// 	PID:   "1222",
	// 	Owner: "test",
	// }

	cmd := exec.Command("ps", "-p", a.PID, "-o", "user,pid,ppid,stat,etime,lstart,comm,args")

	output, err := cmd.Output()

	if err != nil {
		// fmt.Println("Unable to get details")
	}

	rows := strings.Split(string(output), "\n")
	if len(rows) < 2 {
		// fmt.Println("No process details found")
		return
	}

	label := rows[0]
	labels := map[int]string{}
	setLabels(&labels, label)

	portdetails := []PortDetail{}
	setPortDetails(&portdetails, rows[1:], labels)

	if len(portdetails) > 0 {
		a.PortDetails = portdetails[0]
	}
	// fmt.Println("USER:", portdetails[0].User)
	// fmt.Println("PPID:", portdetails[0].PPID)
	// fmt.Println("STAT:", portdetails[0].STAT)
	// fmt.Println("ELAPSED:", portdetails[0].ELAPSED)
	// fmt.Println("STARTED:", portdetails[0].STARTED)
	// fmt.Println("COMMAND:", portdetails[0].COMMAND)
	a.PortDetails = portdetails[0]
}
