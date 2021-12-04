package wsgrp

import (
	"context"
	"github.com/kevguy/algosearch/backend/foundation/web"
	"github.com/kevguy/algosearch/backend/foundation/websocket"
	"net/http"
)

type Handlers struct {
	Hub *websocket.Hub
}

// ServeWS sets up the websocket
func (h Handlers) ServeWS(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	websocket.ServeWs(h.Hub, w, r)
	return nil
}

func (h Handlers) ServeHomePage(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//if r.URL.Path != "/" {
	//	http.Error(w, "Not found", http.StatusNotFound)
	//	return nil
	//}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}
	http.ServeFile(w, r, "home.html")
	return nil
}

func (h Handlers) SendDummy(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var message = "Fuck You"
	h.Hub.ExternalBroadcast <- []byte(message)
	return web.Respond(ctx, w, nil, http.StatusOK)
}
