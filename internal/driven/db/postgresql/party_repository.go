package postgresql

import (
	"context"
	"gitlab.com/ricardo-public/tracing/pkg/tracing"
	"gitlab.com/ricardo134/party-service/internal/core/entities"
	"gitlab.com/ricardo134/party-service/internal/core/ports"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	_, span := tracing.Tracer.Start(ctx, "postgres.partyRepository.Save")
	defer span.End()

	span.SetAttributes(attribute.Int("party.id", int(party.ID)))

	err := p.client.Save(&party).Error
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, notFoundOrElseError(err)
	}

	return &party, err
}

func (p partyRepository) Delete(ctx context.Context, partyID uint) error {
	_, span := tracing.Tracer.Start(ctx, "postgres.partyRepository.Delete")
	defer span.End()

	span.SetAttributes(attribute.Int("party.id", int(partyID)))

	err := p.client.Delete(&entities.Party{}, partyID).Error
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return notFoundOrElseError(err)
	}

	return nil
}

func (p partyRepository) DeleteAllForUser(ctx context.Context, userID uint) error {
	_, span := tracing.Tracer.Start(ctx, "postgres.partyRepository.DeleteAllForUser")
	defer span.End()

	span.SetAttributes(attribute.Int("user.id", int(userID)))

	err := p.client.Where("user_id = ?", userID).Delete(&entities.Party{}).Error
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return notFoundOrElseError(err)
	}

	return nil
}
