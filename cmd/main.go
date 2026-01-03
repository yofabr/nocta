package main

import (
	"nocta/internal/application"
	"nocta/internal/gui"
)

func main() {
	app := application.NewApplication()
	gui.NewGUI(app)
}
