package app

import (
	"context"
	"gitlab.com/ricardo-public/tracing/pkg/tracing"
	"gitlab.com/ricardo134/party-service/internal/core/entities"
	partyPort "gitlab.com/ricardo134/party-service/internal/core/ports"
)

type PartyService interface {
	partyPort.PartyService
}

type partyService struct {
	repo partyPort.PartyRepository
}

func NewPartyService(repo partyPort.PartyRepository) PartyService {
	return partyService{
		repo: repo,
	}
}

func (p partyService) Get(ctx context.Context, partyID uint) (*entities.Party, error) {
	return p.repo.Get(ctx, partyID)
}

func (p partyService) GetAll(ctx context.Context) ([]entities.Party, error) {
	return p.repo.GetAll(ctx)
}

func (p partyService) GetAllForUser(ctx context.Context, userID uint) ([]entities.Party, error) {
	return p.repo.GetAllForUser(ctx, userID)
}

func (p partyService) Save(ctx context.Context, party entities.Party) (*entities.Party, error) {
	nctx, span := tracing.Tracer.Start(ctx, "app.partyService.Save")
	defer span.End()

	return p.repo.Save(nctx, party)
}

func (p partyService) Delete(ctx context.Context, partyID uint) error {
	nctx, span := tracing.Tracer.Start(ctx, "app.partyService.Delete")
	defer span.End()

	return p.repo.Delete(nctx, partyID)
}

func (p partyService) DeleteAllForUser(ctx context.Context, partyID uint) error {
	nctx, span := tracing.Tracer.Start(ctx, "app.partyService.DeleteAllForUser")
	defer span.End()

	return p.repo.Delete(nctx, partyID)
}

func (p partyService) Joined(ctx context.Context, partyID, userID uint) error {
	nctx, span := tracing.Tracer.Start(ctx, "app.partyService.Joined")
	defer span.End()

	return p.repo.Joined(nctx, partyID, userID)
}
