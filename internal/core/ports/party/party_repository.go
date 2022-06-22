package partyPort

import (
	"context"
	"ricardo/party-service/internal/core/entities"
)

type PartyRepository interface {
	Get(ctx context.Context, partyID uint) (*entities.Party, error)
	GetAll(ctx context.Context) ([]entities.Party, error)
	GetAllForUser(ctx context.Context, userID uint) ([]entities.Party, error)
	Save(ctx context.Context, party entities.Party) (*entities.Party, error)
	Delete(ctx context.Context, partyID uint) error
	DeleteAllForUser(ctx context.Context, userID uint) error
}
