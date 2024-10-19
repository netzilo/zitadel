package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/milestone"
)

const (
	MilestonesProjectionTable = "projections.milestones2"

	MilestoneColumnInstanceID  = "instance_id"
	MilestoneColumnType        = "type"
	MilestoneColumnReachedDate = "reached_date"
	MilestoneColumnPushedDate  = "last_pushed_date"
)

type milestoneProjection struct{}

func newMilestoneProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, &milestoneProjection{})
}

func (*milestoneProjection) Name() string {
	return MilestonesProjectionTable
}

func (*milestoneProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(MilestoneColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(MilestoneColumnType, handler.ColumnTypeEnum),
			handler.NewColumn(MilestoneColumnReachedDate, handler.ColumnTypeTimestamp, handler.Nullable()),
			handler.NewColumn(MilestoneColumnPushedDate, handler.ColumnTypeTimestamp, handler.Nullable()),
		},
			handler.NewPrimaryKey(MilestoneColumnInstanceID, MilestoneColumnType),
		),
	)
}

// Reducers implements handler.Projection.
func (p *milestoneProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: milestone.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  milestone.ReachedEventType,
					Reduce: p.reduceReached,
				},
				{
					Event:  milestone.PushedEventType,
					Reduce: p.reducePushed,
				},
			},
		},
	}
}

func (p *milestoneProjection) reduceReached(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*milestone.ReachedEvent](event)
	if err != nil {
		return nil, err
	}
	reached := e.CreatedAt()
	if e.ReachedDate != nil {
		reached = *e.ReachedDate
	}
	return handler.NewCreateStatement(event, []handler.Column{
		handler.NewCol(MilestoneColumnInstanceID, e.Agg.InstanceID),
		handler.NewCol(MilestoneColumnType, e.MilestoneType),
		handler.NewCol(MilestoneColumnReachedDate, reached),
	}), nil
}

func (p *milestoneProjection) reducePushed(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*milestone.PushedEvent](event)
	if err != nil {
		return nil, err
	}
	if e.Agg.Version != milestone.AggregateVersion {
		return handler.NewNoOpStatement(event), nil // Skip v1 events.
	}
	if e.MilestoneType != milestone.InstanceDeleted {
		pushed := e.CreatedAt()
		if e.PushedDate != nil {
			pushed = *e.PushedDate
		}
		return handler.NewUpdateStatement(
			event,
			[]handler.Column{
				handler.NewCol(MilestoneColumnPushedDate, pushed),
			},
			[]handler.Condition{
				handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(MilestoneColumnType, e.MilestoneType),
			},
		), nil
	}
	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}
