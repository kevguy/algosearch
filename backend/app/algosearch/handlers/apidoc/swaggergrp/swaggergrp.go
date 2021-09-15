// Package swaggergrp maintains the group of handlers for serving swagger documentation.
package swaggergrp

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"html/template"
	"io"
	"net/http"
	"os"
)

// Handlers manages the set of swagger endpoints.
type Handlers struct {
	tmpl *template.Template
	swaggerFileName string
	hostProtocol string
	hostEndPoint string
}

func NewIndex(hostProtocol string, hostEndPoint string, fileName string) (Handlers, error) {
	index, err := os.Open("swagger/index.tmpl")
	if err != nil {
		return Handlers{}, errors.Wrap(err, "open index page")
	}
	defer index.Close()
	rawTmpl, err := io.ReadAll(index)
	if err != nil {
		return Handlers{}, errors.Wrap(err, "reading index page")
	}

	tmpl := template.New("index")
	if _, err := tmpl.Parse(string(rawTmpl)); err != nil {
		return Handlers{}, errors.Wrap(err, "creating template")
	}

	fmt.Println(hostEndPoint)
	sg := Handlers{
		tmpl:            tmpl,
		swaggerFileName: fileName,
		hostProtocol: hostProtocol,
		hostEndPoint: hostEndPoint,
	}

	return sg, nil
}

func (h Handlers) ServeDoc(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var markup bytes.Buffer
	vars := map[string]interface{}{
		"HostEndPoint": h.hostEndPoint,
		"HostProtocol": h.hostProtocol,
		//"FileName": "cal-engine-swagger",
		"FileName": h.swaggerFileName,
		//"GraphQLEndpoint": ig.graphQLEndpoint + "/graphql",
		//"MapsKey":         ig.mapsKey,
		//"AuthHeaderName":  ig.authHeaderName,
		//"AuthToken":       ig.authToken,
	}

	fmt.Println(h.hostEndPoint)

	if err := h.tmpl.Execute(&markup, vars); err != nil {
		return errors.Wrap(err, "executing template")
	}

	io.Copy(w, &markup)
	return nil
}
