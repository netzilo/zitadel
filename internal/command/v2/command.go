package command

import "github.com/caos/zitadel/internal/eventstore"

type Command struct {
	es *eventstore.Eventstore
}
