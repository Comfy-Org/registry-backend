package implementation

import (
	"context"
	"registry-backend/drip"
	"registry-backend/ent/user"
	"registry-backend/mapper"

	"github.com/rs/zerolog/log"
)

func (impl *DripStrictServerImplementation) GetUser(ctx context.Context, request drip.GetUserRequestObject) (drip.GetUserResponseObject, error) {
	userId, err := mapper.GetUserIDFromContext(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Stack().Err(err).Msg("")
		return drip.GetUser401Response{}, err
	}

	user, err := impl.Client.User.Query().Where(user.IDEQ(userId)).Only(ctx)

	if (err != nil) || (user == nil) {
		return drip.GetUser404Response{}, err
	}

	return drip.GetUser200JSONResponse{
		Id:         &user.ID,
		Email:      &user.Email,
		Name:       &user.Name,
		IsApproved: &user.IsApproved,
		IsAdmin:    &user.IsAdmin,
	}, nil
}
