// Package samplegrp maintains the group of handlers for sample endpoints.
package samplegrp

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

// Handlers manages the set of sample endpoints.
type Handlers struct {}

func (h Handlers) SendOK(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(status)
	return nil
}

func (h Handlers) SendError(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	//_, ok := ctx.Value(web.KeyValues).(*web.Values)
	//if !ok {
	//	return web.NewShutdownError("web value missing from context")
	//}
	return errors.New( "testing for triggering error")
}
