package main

import (
	"fmt"
	"log"

	"nocta/internal/config"
	"nocta/internal/gui"
	"nocta/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Warning: failed to load config, using defaults: %v", err)
		cfg = &config.Config{}
		cfg.GUI.Width = 900
		cfg.GUI.Height = 520
		cfg.GUI.Split = 0.35
	}

	portService := service.NewPortService()
	gui.NewGUI(portService, cfg)
	fmt.Println("Nocta started")
}
