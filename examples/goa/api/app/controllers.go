// Code generated by goagen v1.3.1, DO NOT EDIT.
//
// API "cellar": Application Controllers
//
// Command:
// $ goagen
// --design=github.com/dikhan/terraform-provider-openapi/examples/goa/api/design
// --out=$(GOPATH)/src/github.com/dikhan/terraform-provider-openapi/examples/goa/api
// --version=v1.3.1

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

// initService sets up the service encoders, decoders and mux.
func initService(service *goa.Service) {
	// Setup encoders and decoders
	service.Encoder.Register(goa.NewJSONEncoder, "application/json")
	service.Encoder.Register(goa.NewGobEncoder, "application/gob", "application/x-gob")
	service.Encoder.Register(goa.NewXMLEncoder, "application/xml")
	service.Decoder.Register(goa.NewJSONDecoder, "application/json")
	service.Decoder.Register(goa.NewGobDecoder, "application/gob", "application/x-gob")
	service.Decoder.Register(goa.NewXMLDecoder, "application/xml")

	// Setup default encoder and decoder
	service.Encoder.Register(goa.NewJSONEncoder, "*/*")
	service.Decoder.Register(goa.NewJSONDecoder, "*/*")
}

// BottleController is the controller interface for the Bottle actions.
type BottleController interface {
	goa.Muxer
	Create(*CreateBottleContext) error
	Show(*ShowBottleContext) error
}

// MountBottleController "mounts" a Bottle resource controller on the given service.
func MountBottleController(service *goa.Service, ctrl BottleController) {
	initService(service)
	var h goa.Handler

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewCreateBottleContext(ctx, req, service)
		if err != nil {
			return err
		}
		// Build the payload
		if rawPayload := goa.ContextRequest(ctx).Payload; rawPayload != nil {
			rctx.Payload = rawPayload.(*BottlePayload)
		} else {
			return goa.MissingPayloadError()
		}
		return ctrl.Create(rctx)
	}
	service.Mux.Handle("POST", "/bottles/", ctrl.MuxHandler("create", h, unmarshalCreateBottlePayload))
	service.LogInfo("mount", "ctrl", "Bottle", "action", "Create", "route", "POST /bottles/")

	h = func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		// Check if there was an error loading the request
		if err := goa.ContextError(ctx); err != nil {
			return err
		}
		// Build the context
		rctx, err := NewShowBottleContext(ctx, req, service)
		if err != nil {
			return err
		}
		return ctrl.Show(rctx)
	}
	service.Mux.Handle("GET", "/bottles/:id", ctrl.MuxHandler("show", h, nil))
	service.LogInfo("mount", "ctrl", "Bottle", "action", "Show", "route", "GET /bottles/:id")
}

// unmarshalCreateBottlePayload unmarshals the request body into the context request data Payload field.
func unmarshalCreateBottlePayload(ctx context.Context, service *goa.Service, req *http.Request) error {
	payload := &bottlePayload{}
	if err := service.DecodeRequest(req, payload); err != nil {
		return err
	}
	if err := payload.Validate(); err != nil {
		// Initialize payload with private data structure so it can be logged
		goa.ContextRequest(ctx).Payload = payload
		return err
	}
	goa.ContextRequest(ctx).Payload = payload.Publicize()
	return nil
}