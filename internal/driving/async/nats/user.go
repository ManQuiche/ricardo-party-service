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

type userHandler struct {
	partyService app.PartyService
	userService  app.UserService
}

type UserHandler interface {
	Created(awt tracing.AnyWithTrace[entities.User])
	Updated(awt tracing.AnyWithTrace[entities.User])
	Deleted(awt tracing.AnyWithTrace[uint])
}

func NewUserHandler(partySvc app.PartyService, userSvc app.UserService) UserHandler {
	return userHandler{partySvc, userSvc}
}

func (nh userHandler) Created(awt tracing.AnyWithTrace[entities.User]) {
	traceID, err := trace.TraceIDFromHex(awt.TraceID)
	if err != nil {
		log.Println(errors.Wrap(err, fmt.Sprintf("cannot parse traceID %s", awt.TraceID)).Error())
	}

	ctx := trace.ContextWithRemoteSpanContext(context.Background(), trace.NewSpanContext(
		trace.SpanContextConfig{
			TraceID: traceID,
		},
	))
	nctx, span := tracing.Tracer.Start(ctx, "nats.UserHandler.UserCreated")
	defer span.End()

	_, _ = nh.userService.Save(nctx, entities.User{
		ID:       awt.Any.ID,
		Username: awt.Any.Username,
	})
}

func (nh userHandler) Updated(awt tracing.AnyWithTrace[entities.User]) {
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

	_, _ = nh.userService.Save(nctx, entities.User{
		ID:       awt.Any.ID,
		Username: awt.Any.Username,
	})
}

func (nh userHandler) Deleted(awt tracing.AnyWithTrace[uint]) {
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
