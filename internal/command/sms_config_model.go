package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

type IAMSMSConfigWriteModel struct {
	eventstore.WriteModel

	ID     string
	Twilio *TwilioConfig
	HTTP   *HTTPConfig
	State  domain.SMSConfigState
}

type TwilioConfig struct {
	Description  string
	SID          string
	Token        *crypto.CryptoValue
	SenderNumber string
}

type HTTPConfig struct {
	Description string
	Endpoint    string
}

func NewIAMSMSConfigWriteModel(instanceID, id string) *IAMSMSConfigWriteModel {
	return &IAMSMSConfigWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   instanceID,
			ResourceOwner: instanceID,
			InstanceID:    instanceID,
		},
		ID: id,
	}
}

func (wm *IAMSMSConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *instance.SMSConfigTwilioAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Twilio = &TwilioConfig{
				SID:          e.SID,
				Token:        e.Token,
				SenderNumber: e.SenderNumber,
			}
			wm.State = domain.SMSConfigStateInactive
		case *instance.SMSConfigTwilioChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			if e.SID != nil {
				wm.Twilio.SID = *e.SID
			}
			if e.SenderNumber != nil {
				wm.Twilio.SenderNumber = *e.SenderNumber
			}
		case *instance.SMSConfigTwilioTokenChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Twilio.Token = e.Token
		case *instance.SMSConfigHTTPAddedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.HTTP = &HTTPConfig{
				Endpoint: e.Endpoint,
			}
			wm.State = domain.SMSConfigStateInactive
		case *instance.SMSConfigHTTPChangedEvent:
			if wm.ID != e.ID {
				continue
			}
			if e.Endpoint != nil {
				wm.HTTP.Endpoint = *e.Endpoint
			}
		case *instance.SMSConfigTwilioActivatedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.SMSConfigStateActive
		case *instance.SMSConfigTwilioDeactivatedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.SMSConfigStateInactive
		case *instance.SMSConfigTwilioRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Twilio = nil
			wm.HTTP = nil
			wm.State = domain.SMSConfigStateRemoved
		case *instance.SMSConfigActivatedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.SMSConfigStateActive
		case *instance.SMSConfigDeactivatedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.State = domain.SMSConfigStateInactive
		case *instance.SMSConfigRemovedEvent:
			if wm.ID != e.ID {
				continue
			}
			wm.Twilio = nil
			wm.HTTP = nil
			wm.State = domain.SMSConfigStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}
func (wm *IAMSMSConfigWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.SMSConfigTwilioAddedEventType,
			instance.SMSConfigTwilioChangedEventType,
			instance.SMSConfigTwilioTokenChangedEventType,
			instance.SMSConfigHTTPAddedEventType,
			instance.SMSConfigHTTPChangedEventType,
			instance.SMSConfigTwilioActivatedEventType,
			instance.SMSConfigTwilioDeactivatedEventType,
			instance.SMSConfigTwilioRemovedEventType,
			instance.SMSConfigActivatedEventType,
			instance.SMSConfigDeactivatedEventType,
			instance.SMSConfigRemovedEventType).
		Builder()
}

func (wm *IAMSMSConfigWriteModel) NewTwilioChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, id string, description, sid, senderNumber *string) (*instance.SMSConfigTwilioChangedEvent, bool, error) {
	changes := make([]instance.SMSConfigTwilioChanges, 0)
	var err error

	if wm.Twilio == nil {
		return nil, false, nil
	}

	if description != nil && wm.Twilio.Description != *description {
		changes = append(changes, instance.ChangeSMSConfigTwilioDescription(*description))
	}
	if sid != nil && wm.Twilio.SID != *sid {
		changes = append(changes, instance.ChangeSMSConfigTwilioSID(*sid))
	}
	if senderNumber != nil && wm.Twilio.SenderNumber != *senderNumber {
		changes = append(changes, instance.ChangeSMSConfigTwilioSenderNumber(*senderNumber))
	}

	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := instance.NewSMSConfigTwilioChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}

func (wm *IAMSMSConfigWriteModel) NewHTTPChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, id string, description, endpoint *string) (*instance.SMSConfigHTTPChangedEvent, bool, error) {
	changes := make([]instance.SMSConfigHTTPChanges, 0)
	var err error

	if wm.HTTP == nil {
		return nil, false, nil
	}

	if description != nil && wm.HTTP.Description != *description {
		changes = append(changes, instance.ChangeSMSConfigHTTPDescription(*description))
	}
	if endpoint != nil && wm.HTTP.Endpoint != *endpoint {
		changes = append(changes, instance.ChangeSMSConfigHTTPEndpoint(*endpoint))
	}

	if len(changes) == 0 {
		return nil, false, nil
	}
	changeEvent, err := instance.NewSMSConfigHTTPChangedEvent(ctx, aggregate, id, changes)
	if err != nil {
		return nil, false, err
	}
	return changeEvent, true, nil
}
