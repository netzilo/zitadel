package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/policy"
)

type InstancePasswordComplexityPolicyWriteModel struct {
	PasswordComplexityPolicyWriteModel
}

func NewInstancePasswordComplexityPolicyWriteModel(instanceID string) *InstancePasswordComplexityPolicyWriteModel {
	return &InstancePasswordComplexityPolicyWriteModel{
		PasswordComplexityPolicyWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
		},
	}
}

func (wm *InstancePasswordComplexityPolicyWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.PasswordComplexityPolicyAddedEvent:
			wm.PasswordComplexityPolicyWriteModel.AppendEvents(&e.PasswordComplexityPolicyAddedEvent)
		case *instance.PasswordComplexityPolicyChangedEvent:
			wm.PasswordComplexityPolicyWriteModel.AppendEvents(&e.PasswordComplexityPolicyChangedEvent)
		}
	}
}

func (wm *InstancePasswordComplexityPolicyWriteModel) Reduce() error {
	return wm.PasswordComplexityPolicyWriteModel.Reduce()
}

func (wm *InstancePasswordComplexityPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.PasswordComplexityPolicyWriteModel.AggregateID).
		EventTypes(
			instance.PasswordComplexityPolicyAddedEventType,
			instance.PasswordComplexityPolicyChangedEventType).
		Builder()
}

func (wm *InstancePasswordComplexityPolicyWriteModel) NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	minLength uint64,
	hasLowercase,
	hasUppercase,
	hasNumber,
	hasSymbol bool,
) (*instance.PasswordComplexityPolicyChangedEvent, bool) {

	changes := make([]policy.PasswordComplexityPolicyChanges, 0)
	if wm.MinLength != minLength {
		changes = append(changes, policy.ChangeMinLength(minLength))
	}
	if wm.HasLowercase != hasLowercase {
		changes = append(changes, policy.ChangeHasLowercase(hasLowercase))
	}
	if wm.HasUppercase != hasUppercase {
		changes = append(changes, policy.ChangeHasUppercase(hasUppercase))
	}
	if wm.HasNumber != hasNumber {
		changes = append(changes, policy.ChangeHasNumber(hasNumber))
	}
	if wm.HasSymbol != hasSymbol {
		changes = append(changes, policy.ChangeHasSymbol(hasSymbol))
	}
	if len(changes) == 0 {
		return nil, false
	}
	changedEvent, err := instance.NewPasswordComplexityPolicyChangedEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, false
	}
	return changedEvent, true
}
