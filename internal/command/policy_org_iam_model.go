package command

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

type PolicyOrgIAMWriteModel struct {
	eventstore.WriteModel

	UserLoginMustBeDomain bool
	State                 domain.PolicyState
}

func (wm *PolicyOrgIAMWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *policy.OrgIAMPolicyAddedEvent:
			wm.UserLoginMustBeDomain = e.UserLoginMustBeDomain
			wm.State = domain.PolicyStateActive
		case *policy.OrgIAMPolicyChangedEvent:
			if e.UserLoginMustBeDomain != nil {
				wm.UserLoginMustBeDomain = *e.UserLoginMustBeDomain
			}
		}
	}
	return wm.WriteModel.Reduce()
}
