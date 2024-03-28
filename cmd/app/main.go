package main

import (
	wsmanager "dummy-chat/internal/manager"
	"dummy-chat/internal/server"
)

func main() {
	manager := wsmanager.New()
	serv := server.New(manager, ":8000")
	go manager.Run()
	serv.Serve()
}
