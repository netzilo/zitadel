package query

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

type IAMMemberReadModel struct {
	MemberReadModel

	userID string
	iamID  string
}

func NewIAMMemberReadModel(iamID, userID string) *IAMMemberReadModel {
	return &IAMMemberReadModel{
		iamID:  iamID,
		userID: userID,
	}
}

func (rm *IAMMemberReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *iam.MemberAddedEvent:
			rm.MemberReadModel.AppendEvents(&e.MemberAddedEvent)
		case *iam.MemberChangedEvent:
			rm.MemberReadModel.AppendEvents(&e.MemberChangedEvent)
		case *member.MemberAddedEvent, *member.MemberChangedEvent, *iam.MemberRemovedEvent:
			rm.MemberReadModel.AppendEvents(e)
		}
	}
}

func (rm *IAMMemberReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, iam.AggregateType).
		AggregateIDs(rm.iamID).
		EventData(map[string]interface{}{
			"userId": rm.userID,
		})
}
