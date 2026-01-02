package application

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Application struct{}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) ListAllPorts() {
	cmd := exec.Command("ss", "-tulnp")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
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

func (a *Application) QueryPort(query QueryParams) {
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

	fmt.Println(string(out))
}
