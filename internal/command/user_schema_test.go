package command

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/user/schema"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_CreateUserSchema(t *testing.T) {
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx        context.Context
		userSchema *CreateUserSchema
	}
	type res struct {
		id      string
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"no type, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:        authz.NewMockContext("instanceID", "", ""),
				userSchema: &CreateUserSchema{},
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-DGFj3", "Errors.UserSchema.Type.Missing"),
			},
		},
		{
			"invalid authenticator, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				userSchema: &CreateUserSchema{
					Type: "type",
					PossibleAuthenticators: []domain.AuthenticatorType{
						domain.AuthenticatorTypeUnspecified,
					},
				},
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-Gh652", "Errors.UserSchema.Authenticator.Invalid"),
			},
		},
		{
			"user schema created",
			fields{
				eventstore: expectEventstore(
					expectPush(
						schema.NewCreatedEvent(
							context.Background(),
							&schema.NewAggregate("id1", "instanceID").Aggregate,
							"type",
							map[string]any{},
							[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
						),
					),
				),
				idGenerator: mock.ExpectID(t, "id1"),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				userSchema: &CreateUserSchema{
					Type: "type",
					PossibleAuthenticators: []domain.AuthenticatorType{
						domain.AuthenticatorTypeUsername,
					},
				},
			},
			res{
				id: "id1",
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore(t),
				idGenerator: tt.fields.idGenerator,
			}
			gotID, gotDetails, err := c.CreateUserSchema(tt.args.ctx, tt.args.userSchema)
			assert.Equal(t, tt.res.id, gotID)
			assert.Equal(t, tt.res.details, gotDetails)
			assert.ErrorIs(t, err, tt.res.err)
		})
	}
}

func TestCommands_UpdateUserSchema(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		userSchema *UpdateUserSchema
	}
	type res struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"missing id, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:        authz.NewMockContext("instanceID", "", ""),
				userSchema: &UpdateUserSchema{},
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-H5421", "Errors.IDMissing"),
			},
		},
		{
			"empty type, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				userSchema: &UpdateUserSchema{
					ID:   "id1",
					Type: gu.Ptr(""),
				},
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-G43gn", "Errors.UserSchema.Type.Missing"),
			},
		},
		{
			"invalid authenticator, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				userSchema: &UpdateUserSchema{
					ID: "id1",
					PossibleAuthenticators: []domain.AuthenticatorType{
						domain.AuthenticatorTypeUnspecified,
					},
				},
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-WF4hg", "Errors.UserSchema.Authenticator.Invalid"),
			},
		},
		{
			"not active / exists, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				userSchema: &UpdateUserSchema{
					ID:   "id1",
					Type: gu.Ptr("type"),
					PossibleAuthenticators: []domain.AuthenticatorType{
						domain.AuthenticatorTypeUsername,
					},
				},
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMA-HB3e1", "Errors.UserSchema.NotActive"),
			},
		},
		{
			"no changes",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								map[string]any{},
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				userSchema: &UpdateUserSchema{
					ID:   "id1",
					Type: gu.Ptr("type"),
					PossibleAuthenticators: []domain.AuthenticatorType{
						domain.AuthenticatorTypeUsername,
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
		{
			"update type",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								map[string]any{},
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
					expectPush(
						schema.NewUpdatedEvent(
							context.Background(),
							&schema.NewAggregate("id1", "instanceID").Aggregate,
							[]schema.Changes{schema.ChangeSchemaType("type", "newType")},
						),
					),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				userSchema: &UpdateUserSchema{
					ID:   "id1",
					Type: gu.Ptr("newType"),
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
		{
			"update possible authenticators",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								map[string]any{},
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
					expectPush(
						schema.NewUpdatedEvent(
							context.Background(),
							&schema.NewAggregate("id1", "instanceID").Aggregate,
							[]schema.Changes{schema.ChangePossibleAuthenticators([]domain.AuthenticatorType{
								domain.AuthenticatorTypeUsername,
								domain.AuthenticatorTypePassword,
							})},
						),
					),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				userSchema: &UpdateUserSchema{
					ID: "id1",
					PossibleAuthenticators: []domain.AuthenticatorType{
						domain.AuthenticatorTypeUsername,
						domain.AuthenticatorTypePassword,
					},
				},
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := c.UpdateUserSchema(tt.args.ctx, tt.args.userSchema)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.details, got)
		})
	}
}

func TestCommands_DeactivateUserSchema(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
		id  string
	}
	type res struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"missing id, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "",
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-Vvf3w", "Errors.IDMissing"),
			},
		},
		{
			"not active / exists, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "id1",
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMA-E4t4z", "Errors.UserSchema.NotActive"),
			},
		},
		{
			"deactivate ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								map[string]any{},
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
					expectPush(
						schema.NewDeactivatedEvent(
							context.Background(),
							&schema.NewAggregate("id1", "instanceID").Aggregate,
						),
					),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "id1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := c.DeactivateUserSchema(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.details, got)
		})
	}
}

func TestCommands_ReactivateUserSchema(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
		id  string
	}
	type res struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"missing id, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "",
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-wq3Gw", "Errors.IDMissing"),
			},
		},
		{
			"not deactivated / exists, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "id1",
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMA-DGzh5", "Errors.UserSchema.NotInactive"),
			},
		},
		{
			"reactivate ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								map[string]any{},
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
						eventFromEventPusher(
							schema.NewDeactivatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
							),
						),
					),
					expectPush(
						schema.NewReactivatedEvent(
							context.Background(),
							&schema.NewAggregate("id1", "instanceID").Aggregate,
						),
					),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "id1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := c.ReactivateUserSchema(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.details, got)
		})
	}
}

func TestCommands_DeleteUserSchema(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
		id  string
	}
	type res struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"missing id, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "",
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-E22gg", "Errors.IDMissing"),
			},
		},
		{
			"not exists, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "id1",
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMA-Grg41", "Errors.UserSchema.NotExists"),
			},
		},
		{
			"delete ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							schema.NewCreatedEvent(
								context.Background(),
								&schema.NewAggregate("id1", "instanceID").Aggregate,
								"type",
								map[string]any{},
								[]domain.AuthenticatorType{domain.AuthenticatorTypeUsername},
							),
						),
					),
					expectPush(
						schema.NewDeletedEvent(
							context.Background(),
							&schema.NewAggregate("id1", "instanceID").Aggregate,
							"type",
						),
					),
				),
			},
			args{
				ctx: authz.NewMockContext("instanceID", "", ""),
				id:  "id1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := c.DeleteUserSchema(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.details, got)
		})
	}
}
