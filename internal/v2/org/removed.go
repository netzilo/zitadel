package org

import (
	"github.com/zitadel/zitadel/v2/internal/v2/eventstore"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
)

const RemovedType = eventTypePrefix + "removed"

type RemovedEvent eventstore.Event[eventstore.EmptyPayload]

var _ eventstore.TypeChecker = (*RemovedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *RemovedEvent) ActionType() string {
	return RemovedType
}

func RemovedEventFromStorage(event *eventstore.StorageEvent) (e *RemovedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-RSPYk", "Errors.Invalid.Event.Type")
	}

	return &RemovedEvent{
		StorageEvent: event,
	}, nil
}
