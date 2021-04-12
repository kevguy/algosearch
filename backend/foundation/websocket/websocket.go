package websocket

import "os"

type App struct {
	hub *Hub
	shutdown chan os.Signal
	//mw []Middleware
}
