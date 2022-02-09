package admin

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	obj_pb "github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/notification/channels/twilio"
	"github.com/caos/zitadel/internal/query"
	admin_pb "github.com/caos/zitadel/pkg/grpc/admin"
	settings_pb "github.com/caos/zitadel/pkg/grpc/settings"
)

func listSMSConfigsToModel(req *admin_pb.ListSMSProviderConfigsRequest) (*query.SMSConfigsSearchQueries, error) {
	offset, limit, asc := object.ListQueryToModel(req.Query)
	return &query.SMSConfigsSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
	}, nil
}

func SMTPConfigToPb(sms *query.SMSConfig) *settings_pb.SMSProviderConfig {
	mapped := &settings_pb.SMSProviderConfig{
		Id:      sms.ID,
		State:   smsStateToPb(sms.State),
		Details: obj_pb.ToViewDetailsPb(sms.Sequence, sms.CreationDate, sms.ChangeDate, sms.AggregateID),
		Config:  SMSConfigToPb(sms),
	}
	return mapped
}

func SMSConfigToPb(app *query.SMSConfig) settings_pb.SMSConfig {
	if app.TwilioConfig != nil {
		return TwilioConfigToPb(app.TwilioConfig)
	}
	return nil
}

func TwilioConfigToPb(twilio *query.Twilio) *settings_pb.SMSProviderConfig_Twilio {
	return &settings_pb.SMSProviderConfig_Twilio{
		Twilio: &settings_pb.TwilioConfig{
			Sid:  twilio.SID,
			From: twilio.From,
		},
	}
}

func smsStateToPb(state domain.SMSConfigState) settings_pb.SMSProviderConfigState {
	switch state {
	case domain.SMSConfigStateInactive:
		return settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_INACTIVE
	case domain.SMSConfigStateActive:
		return settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_ACTIVE
	default:
		return settings_pb.SMSProviderConfigState_SMS_PROVIDER_CONFIG_INACTIVE
	}
}

func AddSMSConfigTwilioToConfig(req *admin_pb.AddSMSProviderConfigTwilioRequest) *twilio.TwilioConfig {
	return &twilio.TwilioConfig{
		SID:   req.Sid,
		From:  req.From,
		Token: req.Token,
	}
}

func UpdateSMSConfigTwilioToConfig(req *admin_pb.UpdateSMSProviderConfigTwilioRequest) *twilio.TwilioConfig {
	return &twilio.TwilioConfig{
		SID:  req.Sid,
		From: req.From,
	}
}
