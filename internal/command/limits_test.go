package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/limits"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestLimits_SetLimits(t *testing.T) {
	type fields func(*testing.T) (*eventstore.Eventstore, id.Generator)
	type args struct {
		ctx           context.Context
		resourceOwner string
		setLimits     *SetLimits
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "create limits, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
						),
					),
					id_mock.NewIDGeneratorExpectIDs(t, "limits1")
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "owner1",
				setLimits: &SetLimits{
					AuditLogRetention: gu.Ptr(time.Hour),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "owner1",
				},
			},
		},
		{
			name: "update limits audit log retention, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Minute)),
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
						),
					),
					nil
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "owner1",
				setLimits: &SetLimits{
					AuditLogRetention: gu.Ptr(time.Hour),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "owner1",
				},
			},
		},
		{
			name: "update limits unblock, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Minute)),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(false)),
								),
							),
						),
					),
					nil
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "owner1",
				setLimits: &SetLimits{
					Block: gu.Ptr(false),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "owner1",
				},
			},
		},
		{
			name: "set limits after resetting limits, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
							eventFromEventPusher(
								limits.NewResetEvent(
									context.Background(),
									&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits2", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
						),
					),
					id_mock.NewIDGeneratorExpectIDs(t, "limits2")
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "owner1",
				setLimits: &SetLimits{
					AuditLogRetention: gu.Ptr(time.Hour),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "owner1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := new(Commands)
			r.eventstore, r.idGenerator = tt.fields(t)
			got, err := r.SetLimits(tt.args.ctx, tt.args.resourceOwner, tt.args.setLimits)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestLimits_SetLimitsBulk(t *testing.T) {
	type fields func(*testing.T) (*eventstore.Eventstore, id.Generator)
	type args struct {
		ctx           context.Context
		setLimitsBulk []*SetLimitsBulk
	}
	type res struct {
		want       *domain.ObjectDetails
		wantTarget []*domain.ObjectDetails
		err        func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "create limits, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
						),
					),
					id_mock.NewIDGeneratorExpectIDs(t, "limits1")
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setLimitsBulk: []*SetLimitsBulk{{
					InstanceID:    "instance1",
					ResourceOwner: "owner1",
					SetLimits: SetLimits{
						AuditLogRetention: gu.Ptr(time.Hour),
					},
				}},
			},
			res: res{
				want: &domain.ObjectDetails{},
				wantTarget: []*domain.ObjectDetails{{
					ResourceOwner: "owner1",
				}},
			},
		},
		{
			name: "update limits audit log retention, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Minute)),
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
						),
					),
					nil
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setLimitsBulk: []*SetLimitsBulk{{
					InstanceID:    "instance1",
					ResourceOwner: "owner1",
					SetLimits: SetLimits{
						AuditLogRetention: gu.Ptr(time.Hour),
					},
				}},
			},
			res: res{
				want: &domain.ObjectDetails{},
				wantTarget: []*domain.ObjectDetails{{
					ResourceOwner: "owner1",
				}},
			},
		},
		{
			name: "update limits unblock, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Minute)),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(false)),
								),
							),
						),
					),
					nil
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setLimitsBulk: []*SetLimitsBulk{{
					InstanceID:    "instance1",
					ResourceOwner: "owner1",
					SetLimits: SetLimits{
						Block: gu.Ptr(false),
					},
				}},
			},
			res: res{
				want: &domain.ObjectDetails{},
				wantTarget: []*domain.ObjectDetails{{
					ResourceOwner: "owner1",
				}},
			},
		},
		{
			name: "set limits after resetting limits, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
							eventFromEventPusher(
								limits.NewResetEvent(
									context.Background(),
									&limits.NewAggregate("limits1", "instance1", "owner1").Aggregate,
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("limits2", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
								),
							),
						),
					),
					id_mock.NewIDGeneratorExpectIDs(t, "limits2")
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instance1"),
				setLimitsBulk: []*SetLimitsBulk{{
					InstanceID:    "instance1",
					ResourceOwner: "owner1",
					SetLimits: SetLimits{
						AuditLogRetention: gu.Ptr(time.Hour),
					},
				},
				},
			},
			res: res{
				want: &domain.ObjectDetails{},
				wantTarget: []*domain.ObjectDetails{{
					ResourceOwner: "owner1",
				}},
			},
		},
		{
			name: "set many limits, ok",
			fields: func(*testing.T) (*eventstore.Eventstore, id.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("blocked-1-1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("unblocked-1-2", "instance1", "owner2").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(false)),
								),
							),
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("blocked-2-1", "instance2", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("blocked-3-1", "instance3", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
							eventFromEventPusher(
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("reset-4-1", "instance4", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
							eventFromEventPusher(
								limits.NewResetEvent(
									context.Background(),
									&limits.NewAggregate("reset-4-1", "instance4", "owner1").Aggregate,
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"instance0",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("create-limits-1", "instance0", "owner0").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("blocked-1-1", "instance1", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(false)),
								),
							),
							eventFromEventPusherWithInstanceID(
								"instance1",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("unblocked-1-2", "instance1", "owner2").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
							eventFromEventPusherWithInstanceID(
								"instance2",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("blocked-2-1", "instance2", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(false)),
								),
							),
							eventFromEventPusherWithInstanceID(
								"instance4",
								limits.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&limits.NewAggregate("reset-4-1", "instance4", "owner1").Aggregate,
										limits.SetEventType,
									),
									limits.ChangeBlock(gu.Ptr(true)),
								),
							),
						),
					),
					id_mock.NewIDGeneratorExpectIDs(t, "create-limits-1", "reset-4-1")
			},
			args: args{
				ctx: context.Background(),
				setLimitsBulk: []*SetLimitsBulk{{
					InstanceID:    "instance0",
					ResourceOwner: "owner0",
					SetLimits: SetLimits{
						Block: gu.Ptr(true),
					},
				}, {
					InstanceID:    "instance1",
					ResourceOwner: "owner1",
					SetLimits: SetLimits{
						Block: gu.Ptr(false),
					},
				}, {
					InstanceID:    "instance1",
					ResourceOwner: "owner2",
					SetLimits: SetLimits{
						Block: gu.Ptr(true),
					},
				}, {
					InstanceID:    "instance2",
					ResourceOwner: "owner1",
					SetLimits: SetLimits{
						Block: gu.Ptr(false),
					},
				}, {
					InstanceID:    "instance3",
					ResourceOwner: "owner1",
					SetLimits: SetLimits{
						Block: gu.Ptr(true),
					},
				}, {
					InstanceID:    "instance4",
					ResourceOwner: "owner1",
					SetLimits: SetLimits{
						Block: gu.Ptr(true),
					},
				}},
			},
			res: res{
				want: &domain.ObjectDetails{},
				wantTarget: []*domain.ObjectDetails{{
					ResourceOwner: "owner0",
				}, {
					ResourceOwner: "owner1",
				}, {
					ResourceOwner: "owner2",
				}, {
					ResourceOwner: "owner1",
				}, {
					ResourceOwner: "owner1",
				}, {
					ResourceOwner: "owner1",
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := new(Commands)
			r.eventstore, r.idGenerator = tt.fields(t)
			gotDetails, gotTargetDetails, err := r.SetLimitsBulk(tt.args.ctx, tt.args.setLimitsBulk)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, gotDetails)
				assert.Equal(t, tt.res.wantTarget, gotTargetDetails)
			}
		})
	}
}

func TestLimits_ResetLimits(t *testing.T) {
	type fields func(*testing.T) *eventstore.Eventstore
	type args struct {
		ctx           context.Context
		resourceOwner string
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "not found",
			fields: func(tt *testing.T) *eventstore.Eventstore {
				return eventstoreExpect(
					tt,
					expectFilter(),
				)
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "instance1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-9JToT", "Errors.Limits.NotFound"))
				},
			},
		},
		{
			name: "already removed",
			fields: func(tt *testing.T) *eventstore.Eventstore {
				return eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							limits.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
									limits.SetEventType,
								),
								limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
							),
						),
						eventFromEventPusher(
							limits.NewResetEvent(context.Background(),
								&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
							),
						),
					),
				)
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "instance1",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "COMMAND-9JToT", "Errors.Limits.NotFound"))
				},
			},
		},
		{
			name: "reset limits, ok",
			fields: func(tt *testing.T) *eventstore.Eventstore {
				return eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							limits.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
									limits.SetEventType,
								),
								limits.ChangeAuditLogRetention(gu.Ptr(time.Hour)),
							),
						),
					),
					expectPush(
						eventFromEventPusherWithInstanceID(
							"instance1",
							limits.NewResetEvent(context.Background(),
								&limits.NewAggregate("limits1", "instance1", "instance1").Aggregate,
							),
						),
					),
				)
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instance1"),
				resourceOwner: "instance1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields(t),
			}
			got, err := r.ResetLimits(tt.args.ctx, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}
