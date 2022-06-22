package party

import (
	"context"
	"ricardo/party-service/internal/core/entities"
	partyPort "ricardo/party-service/internal/core/ports/party"
)

type Service interface {
	partyPort.PartyRepository
}

type partyService struct {
	repo partyPort.PartyRepository
}

func NewPartyService(repo partyPort.PartyRepository) Service {
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
	return p.repo.Save(ctx, party)
}

func (p partyService) Delete(ctx context.Context, partyID uint) error {
	return p.repo.Delete(ctx, partyID)
}

func (p partyService) DeleteAllForUser(ctx context.Context, partyID uint) error {
	return p.repo.Delete(ctx, partyID)
}
