package userservice

import (
	"context"
	"fmt"

	"github.com/aghaghiamh/gocast/QAGame/dto"
	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
)

func (s *Service) GetProfile(ctx context.Context, req dto.UserProfileRequest) (dto.UserProfileResponse, error) {
	const op = "userservice.GetUserProfile"
	user, gErr := s.repo.GetUserByID(ctx, req.UserID)
	if gErr != nil {

		return dto.UserProfileResponse{}, richerr.New(op).
			WithError(gErr).
			WithMessage(fmt.Sprintf("User with %d id does not exist. Please register first.", req.UserID))
	}

	return dto.UserProfileResponse{
		UserInfo: dto.UserInfo{
			UserID:      user.ID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
		},
	}, nil
}
