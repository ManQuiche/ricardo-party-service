package party

import (
	"context"
	"ricardo/party-service/internal/core/entities"
)

type PartyRepository interface {
	Get(ctx context.Context, partyID uint) (*entities.Party, error)
	Save(ctx context.Context, party entities.Party) error
}
