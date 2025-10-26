package main

import (
	server "hospital_management_system/cmd"
	"hospital_management_system/config"
)

func main() {
	// Initialize environment variables and DB
	config.Init()

	// Run server
	server.RunServer()
}
