package models

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type PortDetail struct {
	User    string
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

type PortHistory struct {
	Port      ActivePort
	Timestamp time.Time
	Action    string // "opened", "closed", "changed"
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
