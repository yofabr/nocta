package main

import "nocta/internal/application"

func main() {
	app := application.NewApplication()
	// app.ListAllPorts()
	app.ListAllPorts()
	// app.QueryPort(application.QueryParams{
	// 	Port:     5355,
	// 	Protocol: "tcp",
	// })
}
