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

func NewPartyRepository(client *gorm.DB) partyPort.PartyRepository {
	return partyRepository{client: client}
}

func (p partyRepository) Get(ctx context.Context, partyID uint) (*entities.Party, error) {
	var party entities.Party
	err := p.client.First(&party, partyID).Error
	if err != nil {
		return nil, err
	}

	return &party, err
}

func (p partyRepository) GetAll(ctx context.Context) ([]entities.Party, error) {
	var parties []entities.Party
	err := p.client.Find(&parties).Error
	if err != nil {
		return nil, err
	}

	return parties, nil
}

func (p partyRepository) GetAllForUser(ctx context.Context, userID uint) ([]entities.Party, error) {
	var parties []entities.Party
	err := p.client.Where(&entities.Party{UserID: userID}).Find(&parties).Error
	if err != nil {
		return nil, err
	}

	return parties, nil
}

func (p partyRepository) Save(ctx context.Context, party entities.Party) (*entities.Party, error) {
	err := p.client.Save(&party).Error

	if err != nil {
		return nil, err
	}

	return &party, err
}

func (p partyRepository) Delete(ctx context.Context, partyID uint) error {
	return p.client.Delete(&entities.Party{}, partyID).Error
}
