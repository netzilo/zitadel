package command

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/cache"
	"github.com/zitadel/zitadel/internal/cache/gomap"
	"github.com/zitadel/zitadel/internal/cache/noop"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/milestone"
)

func TestCommands_GetMilestonesReached(t *testing.T) {
	cached := &MilestonesReached{
		InstanceID:                        "cached-id",
		InstanceCreated:                   true,
		AuthenticationSucceededOnInstance: true,
	}

	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	aggregate := milestone.NewAggregate(ctx)

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *MilestonesReached
		wantErr error
	}{
		{
			name: "cached",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "cached-id"),
			},
			want: cached,
		},
		{
			name: "filter error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "no events, all false",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID: "instanceID",
			},
		},
		{
			name: "instance created",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.InstanceCreated, time.Now())),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID:      "instanceID",
				InstanceCreated: true,
			},
		},
		{
			name: "instance auth",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance, time.Now())),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID:                        "instanceID",
				AuthenticationSucceededOnInstance: true,
			},
		},
		{
			name: "project created",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.ProjectCreated, time.Now())),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID:     "instanceID",
				ProjectCreated: true,
			},
		},
		{
			name: "app created",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.ApplicationCreated, time.Now())),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID:         "instanceID",
				ApplicationCreated: true,
			},
		},
		{
			name: "app auth",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnApplication, time.Now())),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID:                           "instanceID",
				AuthenticationSucceededOnApplication: true,
			},
		},
		{
			name: "instance deleted",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.InstanceDeleted, time.Now())),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			want: &MilestonesReached{
				InstanceID:      "instanceID",
				InstanceDeleted: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := gomap.NewCache[milestoneIndex, string, *MilestonesReached](
				context.Background(),
				[]milestoneIndex{milestoneIndexInstanceID},
				cache.CacheConfig{Connector: "memory"},
			)
			cache.Set(context.Background(), cached)

			c := &Commands{
				eventstore: tt.fields.eventstore(t),
				caches: &Caches{
					milestones: cache,
				},
			}
			got, err := c.GetMilestonesReached(tt.args.ctx)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_milestonesCompleted(t *testing.T) {
	c := &Commands{
		caches: &Caches{
			milestones: noop.NewCache[milestoneIndex, string, *MilestonesReached](),
		},
	}
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	arg := &MilestonesReached{
		InstanceID:                           "instanceID",
		InstanceCreated:                      true,
		AuthenticationSucceededOnInstance:    true,
		ProjectCreated:                       true,
		ApplicationCreated:                   true,
		AuthenticationSucceededOnApplication: true,
		InstanceDeleted:                      false,
	}
	c.setCachedMilestonesReached(ctx, arg)
	got, ok := c.getCachedMilestonesReached(ctx)
	assert.True(t, ok)
	assert.Equal(t, arg, got)
}

func TestCommands_MilestonePushed(t *testing.T) {
	aggregate := milestone.NewInstanceAggregate("instanceID")
	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	now := time.Now()
	type args struct {
		ctx           context.Context
		instanceID    string
		msType        milestone.Type
		endpoints     []string
		primaryDomain string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "milestone pushed",
			fields: fields{
				eventstore: expectEventstore(
					expectPush(
						milestone.NewPushedEvent(
							context.Background(),
							aggregate,
							milestone.ApplicationCreated,
							now,
							[]string{"foo.com", "bar.com"},
							"example.com",
							"zitadel.com",
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				instanceID:    "instanceID",
				msType:        milestone.ApplicationCreated,
				endpoints:     []string{"foo.com", "bar.com"},
				primaryDomain: "zitadel.com",
			},
			wantErr: nil,
		},
		{
			name: "pusher error",
			fields: fields{
				eventstore: expectEventstore(
					expectPushFailed(
						io.ErrClosedPipe,
						milestone.NewPushedEvent(
							context.Background(),
							aggregate,
							milestone.ApplicationCreated,
							now,
							[]string{"foo.com", "bar.com"},
							"example.com",
							"zitadel.com",
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				instanceID:    "instanceID",
				msType:        milestone.ApplicationCreated,
				endpoints:     []string{"foo.com", "bar.com"},
				primaryDomain: "zitadel.com",
			},
			wantErr: io.ErrClosedPipe,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:     tt.fields.eventstore(t),
				externalDomain: "example.com",
			}
			err := c.MilestonePushed(tt.args.ctx, tt.args.instanceID, tt.args.msType, now, tt.args.endpoints, tt.args.primaryDomain)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommands_oidcSessionMilestones(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	ctx = authz.WithConsoleClientID(ctx, "console")
	now := time.Now()
	aggregate := milestone.NewAggregate(ctx)

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx      context.Context
		clientID string
		isHuman  bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "get error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "client",
				isHuman:  true,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "milestones already reached",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance, now)),
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnApplication, now)),
					),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "client",
				isHuman:  true,
			},
			wantErr: nil,
		},
		{
			name: "auth on instance",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance, now),
					),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "console",
				isHuman:  true,
			},
			wantErr: nil,
		},
		{
			name: "subsequent console login, no milestone",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance, now)),
					),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "console",
				isHuman:  true,
			},
			wantErr: nil,
		},
		{
			name: "subsequent machine login, no milestone",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance, now)),
					),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "client",
				isHuman:  false,
			},
			wantErr: nil,
		},
		{
			name: "auth on app",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance, now)),
					),
					expectPush(
						milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnApplication, now),
					),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "client",
				isHuman:  true,
			},
			wantErr: nil,
		},
		{
			name: "pusher error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnInstance, now)),
					),
					expectPushFailed(
						io.ErrClosedPipe,
						milestone.NewReachedEvent(ctx, aggregate, milestone.AuthenticationSucceededOnApplication, now),
					),
				),
			},
			args: args{
				ctx:      ctx,
				clientID: "client",
				isHuman:  true,
			},
			wantErr: io.ErrClosedPipe,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
				caches: &Caches{
					milestones: noop.NewCache[milestoneIndex, string, *MilestonesReached](),
				},
			}
			err := c.oidcSessionMilestones(tt.args.ctx, tt.args.clientID, tt.args.isHuman, now)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommands_projectCreatedMilestone(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	systemCtx := authz.SetCtxData(ctx, authz.CtxData{
		SystemMemberships: authz.Memberships{
			&authz.Membership{
				MemberType: authz.MemberTypeSystem,
			},
		},
	})
	now := time.Now()
	aggregate := milestone.NewAggregate(ctx)

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "system user",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: systemCtx,
			},
			wantErr: nil,
		},
		{
			name: "get error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "milestone already reached",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.ProjectCreated, now)),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: nil,
		},
		{
			name: "milestone pushed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						milestone.NewReachedEvent(ctx, aggregate, milestone.ProjectCreated, now),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: nil,
		},
		{
			name: "pusher error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPushFailed(
						io.ErrClosedPipe,
						milestone.NewReachedEvent(ctx, aggregate, milestone.ProjectCreated, now),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: io.ErrClosedPipe,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
				caches: &Caches{
					milestones: noop.NewCache[milestoneIndex, string, *MilestonesReached](),
				},
			}
			err := c.projectCreatedMilestone(tt.args.ctx, now)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommands_applicationCreatedMilestone(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	systemCtx := authz.SetCtxData(ctx, authz.CtxData{
		SystemMemberships: authz.Memberships{
			&authz.Membership{
				MemberType: authz.MemberTypeSystem,
			},
		},
	})
	now := time.Now()
	aggregate := milestone.NewAggregate(ctx)

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "system user",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				ctx: systemCtx,
			},
			wantErr: nil,
		},
		{
			name: "get error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "milestone already reached",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(milestone.NewReachedEvent(ctx, aggregate, milestone.ApplicationCreated, now)),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: nil,
		},
		{
			name: "milestone pushed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						milestone.NewReachedEvent(ctx, aggregate, milestone.ApplicationCreated, now),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: nil,
		},
		{
			name: "pusher error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPushFailed(
						io.ErrClosedPipe,
						milestone.NewReachedEvent(ctx, aggregate, milestone.ApplicationCreated, now),
					),
				),
			},
			args: args{
				ctx: ctx,
			},
			wantErr: io.ErrClosedPipe,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
				caches: &Caches{
					milestones: noop.NewCache[milestoneIndex, string, *MilestonesReached](),
				},
			}
			err := c.applicationCreatedMilestone(tt.args.ctx, now)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommands_instanceRemovedMilestone(t *testing.T) {
	ctx := authz.WithInstanceID(context.Background(), "instanceID")
	now := time.Now()
	aggregate := milestone.NewInstanceAggregate("instanceID")

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "milestone pushed",
			fields: fields{
				eventstore: expectEventstore(
					expectPush(
						milestone.NewReachedEvent(ctx, aggregate, milestone.InstanceDeleted, now),
					),
				),
			},
			args: args{
				ctx:        ctx,
				instanceID: "instanceID",
			},
			wantErr: nil,
		},
		{
			name: "pusher error",
			fields: fields{
				eventstore: expectEventstore(
					expectPushFailed(
						io.ErrClosedPipe,
						milestone.NewReachedEvent(ctx, aggregate, milestone.InstanceDeleted, now),
					),
				),
			},
			args: args{
				ctx:        ctx,
				instanceID: "instanceID",
			},
			wantErr: io.ErrClosedPipe,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
				caches: &Caches{
					milestones: noop.NewCache[milestoneIndex, string, *MilestonesReached](),
				},
			}
			err := c.instanceRemovedMilestone(tt.args.ctx, tt.args.instanceID, now)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func (c *Commands) setMilestonesCompletedForTest(instanceID string) {
	c.milestonesCompleted.Store(instanceID, struct{}{})
}
