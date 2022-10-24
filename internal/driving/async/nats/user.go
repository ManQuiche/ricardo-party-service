package nats

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/ricardo-public/tracing/pkg/tracing"
	"gitlab.com/ricardo134/party-service/internal/core/app"
	"gitlab.com/ricardo134/party-service/internal/core/entities"
	"go.opentelemetry.io/otel/trace"
	"log"
)

type handler struct {
	partyService app.PartyService
	userService  app.UserService
}

type Handler interface {
	Created(awt tracing.AnyWithTrace[entities.User])
	Updated(awt tracing.AnyWithTrace[entities.User])
	Deleted(awt tracing.AnyWithTrace[uint])
}

func NewUserHandler(partySvc app.PartyService, userSvc app.UserService) Handler {
	return handler{partySvc, userSvc}
}

func (nh handler) Created(awt tracing.AnyWithTrace[entities.User]) {
	traceID, err := trace.TraceIDFromHex(awt.TraceID)
	if err != nil {
		log.Println(errors.Wrap(err, fmt.Sprintf("cannot parse traceID %s", awt.TraceID)).Error())
	}

	ctx := trace.ContextWithRemoteSpanContext(context.Background(), trace.NewSpanContext(
		trace.SpanContextConfig{
			TraceID: traceID,
		},
	))
	nctx, span := tracing.Tracer.Start(ctx, "nats.UserHandler.Created")
	defer span.End()

	_, _ = nh.userService.Save(nctx, awt.Any)
}

func (nh handler) Updated(awt tracing.AnyWithTrace[entities.User]) {
	traceID, err := trace.TraceIDFromHex(awt.TraceID)
	if err != nil {
		log.Println(errors.Wrap(err, fmt.Sprintf("cannot parse traceID %s", awt.TraceID)).Error())
	}

	ctx := trace.ContextWithRemoteSpanContext(context.Background(), trace.NewSpanContext(
		trace.SpanContextConfig{
			TraceID: traceID,
		},
	))
	nctx, span := tracing.Tracer.Start(ctx, "nats.UserHandler.Updated")
	defer span.End()

	_, _ = nh.userService.Save(nctx, awt.Any)
}

func (nh handler) Deleted(awt tracing.AnyWithTrace[uint]) {
	traceID, err := trace.TraceIDFromHex(awt.TraceID)
	if err != nil {
		log.Println(errors.Wrap(err, fmt.Sprintf("cannot parse traceID %s", awt.TraceID)).Error())
	}

	ctx := trace.ContextWithRemoteSpanContext(context.Background(), trace.NewSpanContext(
		trace.SpanContextConfig{
			TraceID: traceID,
		},
	))
	nctx, span := tracing.Tracer.Start(ctx, "nats.UserHandler.Deleted")
	defer span.End()

	_ = nh.partyService.DeleteAllForUser(nctx, awt.Any)
	_ = nh.userService.Delete(nctx, awt.Any)
}
