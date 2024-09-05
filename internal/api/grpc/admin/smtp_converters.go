package admin

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	admin_pb "github.com/zitadel/zitadel/pkg/grpc/admin"
	settings_pb "github.com/zitadel/zitadel/pkg/grpc/settings"
)

func listSMTPConfigsToModel(req *admin_pb.ListSMTPConfigsRequest) (*query.SMTPConfigsSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &query.SMTPConfigsSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
	}, nil
}

func SMTPConfigToProviderPb(config *query.SMTPConfig) *settings_pb.SMTPConfig {
	return &settings_pb.SMTPConfig{
		Details:       object.ToViewDetailsPb(config.Sequence, config.CreationDate, config.ChangeDate, config.ResourceOwner),
		Id:            config.ID,
		Description:   config.Description,
		Tls:           config.SMTPConfig.TLS,
		Host:          config.SMTPConfig.Host,
		User:          config.SMTPConfig.User,
		State:         SMTPConfigStateToPb(config.State),
		SenderAddress: config.SMTPConfig.SenderAddress,
		SenderName:    config.SMTPConfig.SenderName,
	}
}

func SMTPConfigsToPb(configs []*query.SMTPConfig) []*settings_pb.SMTPConfig {
	c := make([]*settings_pb.SMTPConfig, len(configs))
	for i, config := range configs {
		c[i] = SMTPConfigToProviderPb(config)
	}
	return c
}

func SMTPConfigStateToPb(state domain.SMTPConfigState) settings_pb.SMTPConfigState {
	switch state {
	case domain.SMTPConfigStateUnspecified, domain.SMTPConfigStateRemoved:
		return settings_pb.SMTPConfigState_SMTP_CONFIG_STATE_UNSPECIFIED
	case domain.SMTPConfigStateActive:
		return settings_pb.SMTPConfigState_SMTP_CONFIG_ACTIVE
	case domain.SMTPConfigStateInactive:
		return settings_pb.SMTPConfigState_SMTP_CONFIG_INACTIVE
	default:
		return settings_pb.SMTPConfigState_SMTP_CONFIG_STATE_UNSPECIFIED
	}
}
