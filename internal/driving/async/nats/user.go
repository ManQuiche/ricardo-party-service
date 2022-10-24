package nats

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/ricardo-public/tracing/pkg/tracing"
	"gitlab.com/ricardo134/party-service/internal/core/app"
	"gitlab.com/ricardo134/party-service/internal/core/entities"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log"
)

type handler struct {
	partyService app.PartyService
	userService  app.UserService
}

type Handler interface {
	Created(awt tracing.AnyWithTrace)
	Updated(awt tracing.AnyWithTrace)
	Deleted(awt tracing.AnyWithTrace)
}

func NewUserHandler(partySvc app.PartyService, userSvc app.UserService) Handler {
	return handler{partySvc, userSvc}
}

func (nh handler) Created(awt tracing.AnyWithTrace) {
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

	user, ok := awt.Any.(entities.User)
	if ok == false {
		err = errors.New("could not decode incoming user")
		log.Println(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	_, _ = nh.userService.Save(nctx, user)
}

func (nh handler) Updated(awt tracing.AnyWithTrace) {
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

	user, ok := awt.Any.(entities.User)
	if ok == false {
		err = errors.New("could not decode incoming user")
		log.Println(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	_, _ = nh.userService.Save(nctx, user)
}

func (nh handler) Deleted(awt tracing.AnyWithTrace) {
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

	userID, ok := awt.Any.(uint)
	if ok == false {
		err = errors.New("could not decode incoming user")
		log.Println(err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	_ = nh.partyService.DeleteAllForUser(nctx, userID)
	_ = nh.userService.Delete(nctx, userID)
}
