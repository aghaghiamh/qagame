package presenceservice

import (
	"context"
	"fmt"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
)

func (s Service) Upsert(ctx context.Context, req dto.PresenceUpsertRequest) (dto.PresenceUpsertResponse, error) {
	const op = richerr.Operation("presenceservice.Upsert")

	key := fmt.Sprintf("%s:%d", s.config.Prefix, req.UserID)
	if err := s.repo.Upsert(ctx, key, req.Timestamp, s.config.ExpectedOnlineTime); err != nil {
		return dto.PresenceUpsertResponse{}, richerr.New(op).WithError(err)
	}

	// TODO: proper response
	return dto.PresenceUpsertResponse{}, nil
}
