package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Respond converts a Go value to JSON and sends it to the client.
func Respond(ctx context.Context, w http.ResponseWriter, data interface{}, statusCode int) error {
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "foundation.web.respond")
	span.SetAttributes(attribute.Int("statusCode", statusCode))
	defer span.End()

	// Set the status code for the request logger middleware.
	SetStatusCode(ctx, statusCode)

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	fmt.Println("hihi")

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}

func Redirect(ctx context.Context, w http.ResponseWriter, r *http.Request, url string) error {
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "foundation.web.respond")
	span.SetAttributes(attribute.Int("statusCode", http.StatusMovedPermanently))
	defer span.End()

	// Set the status code for the request logger middleware.
	// If the context is missing this value, don't set it and
	// make sure a response is provided.
	SetStatusCode(ctx, http.StatusMovedPermanently)

	// Write the status code to the response.
	w.WriteHeader(http.StatusMovedPermanently)

	// Send the result back to the client.
	http.Redirect(w, r, url, http.StatusMovedPermanently)
	return nil
}

// Respond converts a Go value to JSON and sends it to the client.
func RespondStr(ctx context.Context, w http.ResponseWriter, data string, statusCode int) error {
	ctx, span := otel.GetTracerProvider().Tracer("").Start(ctx, "foundation.web.respondStr")
	span.SetAttributes(attribute.Int("statusCode", statusCode))
	defer span.End()

	// Set the status code for the request logger middleware.
	// If the context is missing this value, don't set it and
	// make sure a response is provided.
	SetStatusCode(ctx, statusCode)

	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Set the content type and headers once we know marshaling has succeeded.
	//w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}
