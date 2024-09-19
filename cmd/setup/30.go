package setup

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/query/projection"
	"github.com/zitadel/zitadel/v2/internal/repository/instance"
)

type FillFieldsForOrgDomainVerified struct {
	eventstore *eventstore.Eventstore
}

func (mig *FillFieldsForOrgDomainVerified) Execute(ctx context.Context, _ eventstore.Event) error {
	instances, err := mig.eventstore.InstanceIDs(
		ctx,
		0,
		true,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).
			OrderDesc().
			AddQuery().
			AggregateTypes("instance").
			EventTypes(instance.InstanceAddedEventType).
			Builder(),
	)
	if err != nil {
		return err
	}
	for _, instance := range instances {
		ctx := authz.WithInstanceID(ctx, instance)
		if err := projection.OrgDomainVerifiedFields.Trigger(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (mig *FillFieldsForOrgDomainVerified) String() string {
	return "30_fill_fields_for_org_domain_verified"
}
