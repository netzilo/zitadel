package user

import (
	"time"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type humanInitialCodeAddedPayload struct {
	Expiry            time.Duration `json:"expiry,omitempty"`
	TriggeredAtOrigin string        `json:"triggerOrigin,omitempty"`
}

type HumanInitialCodeAddedEvent humanInitialCodeAddedEvent
type humanInitialCodeAddedEvent = eventstore.StorageEvent[humanInitialCodeAddedPayload]

func HumanInitialCodeAddedEventFromStorage(e *eventstore.StorageEvent[eventstore.StoragePayload]) (*HumanInitialCodeAddedEvent, error) {
	event, err := eventstore.EventFromStorage[humanInitialCodeAddedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanInitialCodeAddedEvent)(event), nil
}
