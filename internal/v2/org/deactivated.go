package org

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const DeactivatedType = "org.deactivated"

type DeactivatedEvent eventstore.Event[eventstore.EmptyPayload]

var _ eventstore.TypeChecker = (*DeactivatedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DeactivatedEvent) ActionType() string {
	return DeactivatedType
}

func DeactivatedEventFromStorage(event *eventstore.StorageEvent) (e *DeactivatedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-jeeON", "Errors.Invalid.Event.Type")
	}

	return &DeactivatedEvent{
		StorageEvent: event,
	}, nil
}
