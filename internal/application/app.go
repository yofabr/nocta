package application

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type ActivePort struct {
	Protocol       string
	State          string
	Addr           string
	Port           string
	RecvQ          string
	SendQ          string
	Peer_Addr_Port string
	Process        string
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

func (a *Application) setActivePorts(active_ports []ActivePort) {
	a.ActivePorts = active_ports
}

func setLabels(labels *map[int]string, label string) {
	for i, l := range strings.Split(strings.TrimSpace(label), " ") {
		if l == "" {
			continue
		}
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
		Peer_Addr_Port: port[5],
	}

	return active_port
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

func formatPortString(port_string string) []ActivePort {
	var activePorts []ActivePort

	cmd_list := strings.Split(port_string, "\n")

	label := cmd_list[0]
	labels := map[int]string{}
	setLabels(&labels, label)

	setActivePorts(&activePorts, cmd_list[1:])

	for _, n := range activePorts {
		fmt.Println(n)
	}
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
