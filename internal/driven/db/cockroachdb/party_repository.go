package cockroachdb

import (
	"context"
	"gorm.io/gorm"
	"ricardo/party-service/internal/core/entities"
	"ricardo/party-service/internal/core/ports/party"
)

type partyRepository struct {
	client *gorm.DB
}

func NewPartyRepository(client *gorm.DB) party.PartyRepository {
	return partyRepository{client: client}
}

func (p partyRepository) Get(ctx context.Context, partyID uint) (*entities.Party, error) {
	//TODO implement me
	panic("implement me")
}

func (p partyRepository) Save(ctx context.Context, party entities.Party) error {
	//TODO implement me
	panic("implement me")
}
