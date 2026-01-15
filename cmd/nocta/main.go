package main

import (
	"nocta/internal/gui"
	"nocta/internal/service"
)

func main() {
	portService := service.NewPortService()
	gui.NewGUI(portService)
}
