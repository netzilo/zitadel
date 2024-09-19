package user

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/v2/internal/api/grpc/object"
	"github.com/zitadel/zitadel/v2/internal/user/model"
	"github.com/zitadel/zitadel/v2/pkg/grpc/user"
)

func RefreshTokensToPb(refreshTokens []*model.RefreshTokenView) []*user.RefreshToken {
	tokens := make([]*user.RefreshToken, len(refreshTokens))
	for i, token := range refreshTokens {
		tokens[i] = RefreshTokenToPb(token)
	}
	return tokens
}

func RefreshTokenToPb(token *model.RefreshTokenView) *user.RefreshToken {
	return &user.RefreshToken{
		Id:             token.ID,
		Details:        object.ToViewDetailsPb(token.Sequence, token.CreationDate, token.ChangeDate, token.ResourceOwner),
		ClientId:       token.ClientID,
		AuthTime:       timestamppb.New(token.AuthTime),
		IdleExpiration: timestamppb.New(token.IdleExpiration),
		Expiration:     timestamppb.New(token.Expiration),
		Scopes:         token.Scopes,
		Audience:       token.Audience,
	}
}
