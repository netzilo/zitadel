package command

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user/schemauser"
)

type UserV3WriteModel struct {
	eventstore.WriteModel

	PhoneWM bool
	EmailWM bool
	DataWM  bool

	SchemaID       string
	SchemaRevision uint64

	Email                    string
	IsEmailVerified          bool
	EmailVerifiedFailedCount int
	Phone                    string
	IsPhoneVerified          bool
	PhoneVerifiedFailedCount int

	Data json.RawMessage

	State domain.UserState
}

func NewExistsUserV3WriteModel(resourceOwner, userID string) *UserV3WriteModel {
	return &UserV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		PhoneWM: false,
		EmailWM: false,
		DataWM:  false,
	}
}

func NewUserV3WriteModel(resourceOwner, userID string) *UserV3WriteModel {
	return &UserV3WriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
		PhoneWM: true,
		EmailWM: true,
		DataWM:  true,
	}
}

func (wm *UserV3WriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *schemauser.CreatedEvent:
			wm.SchemaID = e.SchemaID
			wm.SchemaRevision = 1
			wm.Data = e.Data

			wm.State = domain.UserStateActive
		case *schemauser.UpdatedEvent:
			if e.SchemaID != nil {
				wm.SchemaID = *e.SchemaID
			}
			if e.SchemaRevision != nil {
				wm.SchemaRevision = *e.SchemaRevision
			}
			if len(e.Data) > 0 {
				wm.Data = e.Data
			}
		case *schemauser.DeletedEvent:
			wm.State = domain.UserStateDeleted
		case *schemauser.EmailUpdatedEvent:
			wm.Email = string(e.EmailAddress)
		case *schemauser.EmailCodeAddedEvent:
			wm.IsEmailVerified = false
			wm.EmailVerifiedFailedCount = 0
		case *schemauser.EmailVerifiedEvent:
			wm.IsEmailVerified = true
			wm.EmailVerifiedFailedCount = 0
		case *schemauser.EmailVerificationFailedEvent:
			wm.EmailVerifiedFailedCount += 1
		case *schemauser.PhoneChangedEvent:
			wm.Phone = string(e.PhoneNumber)
		case *schemauser.PhoneCodeAddedEvent:
			wm.IsPhoneVerified = false
			wm.PhoneVerifiedFailedCount = 0
		case *schemauser.PhoneVerifiedEvent:
			wm.PhoneVerifiedFailedCount = 0
			wm.IsPhoneVerified = true
		case *schemauser.PhoneVerificationFailedEvent:
			wm.PhoneVerifiedFailedCount += 1
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *UserV3WriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(schemauser.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			schemauser.CreatedType,
			schemauser.DeletedType,
		)
	if wm.DataWM {
		query = query.EventTypes(
			schemauser.UpdatedType,
		)
	}
	if wm.EmailWM {
		query = query.EventTypes(
			schemauser.EmailUpdatedType,
			schemauser.EmailVerifiedType,
			schemauser.EmailCodeAddedType,
			schemauser.EmailVerificationFailedType,
		)
	}
	if wm.PhoneWM {
		query = query.EventTypes(
			schemauser.PhoneUpdatedType,
			schemauser.PhoneVerifiedType,
			schemauser.PhoneCodeAddedType,
			schemauser.PhoneVerificationFailedType,
		)
	}
	return query.Builder()
}

func (wm *UserV3WriteModel) NewUpdatedEvent(
	ctx context.Context,
	agg *eventstore.Aggregate,
	schemaID *string,
	schemaRevision *uint64,
	data json.RawMessage,
) *schemauser.UpdatedEvent {
	changes := make([]schemauser.Changes, 0)
	if schemaID != nil && wm.SchemaID != *schemaID {
		changes = append(changes, schemauser.ChangeSchemaID(wm.SchemaID, *schemaID))
	}
	if schemaRevision != nil && wm.SchemaRevision != *schemaRevision {
		changes = append(changes, schemauser.ChangeSchemaRevision(wm.SchemaRevision, *schemaRevision))
	}
	if !bytes.Equal(wm.Data, data) {
		changes = append(changes, schemauser.ChangeData(data))
	}
	if len(changes) == 0 {
		return nil
	}
	return schemauser.NewUpdatedEvent(ctx, agg, changes)
}

func UserV3AggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            wm.AggregateID,
		Type:          schemauser.AggregateType,
		ResourceOwner: wm.ResourceOwner,
		InstanceID:    wm.InstanceID,
		Version:       schemauser.AggregateVersion,
	}
}

func (wm *UserV3WriteModel) Exists() bool {
	return wm.State != domain.UserStateDeleted && wm.State != domain.UserStateUnspecified
}
