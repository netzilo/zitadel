package project

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

func init() {
	eventstore.RegisterFilterEventMapper(AggregateType, ProjectAddedType, ProjectAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ProjectChangedType, ProjectChangeEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ProjectDeactivatedType, ProjectDeactivatedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ProjectReactivatedType, ProjectReactivatedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ProjectRemovedType, ProjectRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, MemberAddedType, MemberAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, MemberChangedType, MemberChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, MemberRemovedType, MemberRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, MemberCascadeRemovedType, MemberCascadeRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, RoleAddedType, RoleAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, RoleChangedType, RoleChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, RoleRemovedType, RoleRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GrantAddedType, GrantAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GrantChangedType, GrantChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GrantCascadeChangedType, GrantCascadeChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GrantDeactivatedType, GrantDeactivateEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GrantReactivatedType, GrantReactivatedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GrantRemovedType, GrantRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GrantMemberAddedType, GrantMemberAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GrantMemberChangedType, GrantMemberChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GrantMemberRemovedType, GrantMemberRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, GrantMemberCascadeRemovedType, GrantMemberCascadeRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ApplicationAddedType, ApplicationAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ApplicationChangedType, ApplicationChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ApplicationRemovedType, ApplicationRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ApplicationDeactivatedType, ApplicationDeactivatedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ApplicationReactivatedType, ApplicationReactivatedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, OIDCConfigAddedType, OIDCConfigAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, OIDCConfigChangedType, OIDCConfigChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, OIDCConfigSecretChangedType, OIDCConfigSecretChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, OIDCClientSecretCheckSucceededType, OIDCConfigSecretCheckSucceededEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, OIDCClientSecretCheckFailedType, OIDCConfigSecretCheckFailedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, APIConfigAddedType, APIConfigAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, APIConfigChangedType, APIConfigChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, APIConfigSecretChangedType, APIConfigSecretChangedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ApplicationKeyAddedEventType, ApplicationKeyAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, ApplicationKeyRemovedEventType, ApplicationKeyRemovedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SAMLConfigAddedType, SAMLConfigAddedEventMapper)
	eventstore.RegisterFilterEventMapper(AggregateType, SAMLConfigChangedType, SAMLConfigChangedEventMapper)
}
