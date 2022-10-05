package postgresql

import (
	"context"
	"gitlab.com/ricardo134/party-service/internal/core/entities"
	"gitlab.com/ricardo134/party-service/internal/core/ports"
	"gorm.io/gorm"
)

type partyRepository struct {
	client *gorm.DB
}

func NewPartyRepository(client *gorm.DB) ports.PartyRepository {
	return partyRepository{client: client}
}

func (p partyRepository) Get(ctx context.Context, partyID uint) (*entities.Party, error) {
	var party entities.Party
	err := p.client.First(&party, partyID).Error
	if err != nil {
		return nil, notFoundOrElseError(err)
	}

	return &party, nil
}

func (p partyRepository) GetAll(ctx context.Context) ([]entities.Party, error) {
	var parties []entities.Party
	err := p.client.Find(&parties).Error
	if err != nil {
		return nil, notFoundOrElseError(err)
	}

	return parties, nil
}

func (p partyRepository) GetAllForUser(ctx context.Context, userID uint) ([]entities.Party, error) {
	var parties []entities.Party
	err := p.client.Where("user_id = ?", userID).Find(&parties).Error
	if err != nil {
		return nil, notFoundOrElseError(err)
	}

	return parties, nil
}

func (p partyRepository) Save(ctx context.Context, party entities.Party) (*entities.Party, error) {
	err := p.client.Save(&party).Error

	if err != nil {
		return nil, notFoundOrElseError(err)
	}

	return &party, err
}

func (p partyRepository) Delete(ctx context.Context, partyID uint) error {
	err := p.client.Delete(&entities.Party{}, partyID).Error
	if err != nil {
		return notFoundOrElseError(err)
	}

	return nil
}

func (p partyRepository) DeleteAllForUser(ctx context.Context, userID uint) error {
	err := p.client.Where("user_id = ?", userID).Delete(&entities.Party{}).Error
	if err != nil {
		return notFoundOrElseError(err)
	}

	return nil
}
