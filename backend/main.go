package main

import "github.com/google/uuid"

type Lobbies = map[uuid.UUID]*Lobby

func main() {
	app := Application{
		Lobbies: Lobbies{},
	}
	app.ParseFlags()
	app.setupLogging()
	app.SetupHandlers()
	app.Run()
}
