package milestone

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

//go:generate enumer -type Type -json -linecomment -transform=snake
type Type int

const (
	InstanceCreated Type = iota
	AuthenticationSucceededOnInstance
	ProjectCreated
	ApplicationCreated
	AuthenticationSucceededOnApplication
	InstanceDeleted
)

const (
	eventTypePrefix  = "milestone."
	ReachedEventType = eventTypePrefix + "reached"
	PushedEventType  = eventTypePrefix + "pushed"
)

type ReachedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	MilestoneType         Type       `json:"type"`
	ReachedDate           *time.Time `json:"reachedDate,omitempty"` // Defaults to [eventstore.BaseEvent.Creation] when empty
}

// Payload implements eventstore.Command.
func (e *ReachedEvent) Payload() any {
	return e
}

func (e *ReachedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *ReachedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func NewReachedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	typ Type,
) *ReachedEvent {
	return NewReachedEventWithDate(ctx, aggregate, typ, nil)
}

// NewReachedEventWithDate creates a [ReachedEvent] with a fixed Reached Date.
func NewReachedEventWithDate(
	ctx context.Context,
	aggregate *Aggregate,
	typ Type,
	reachedDate *time.Time,
) *ReachedEvent {
	return &ReachedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			&aggregate.Aggregate,
			ReachedEventType,
		),
		MilestoneType: typ,
		ReachedDate:   reachedDate,
	}
}

type PushedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	MilestoneType         Type       `json:"type"`
	ExternalDomain        string     `json:"externalDomain"`
	PrimaryDomain         string     `json:"primaryDomain"`
	Endpoints             []string   `json:"endpoints"`
	PushedDate            *time.Time `json:"pushedDate,omitempty"` // Defaults to [eventstore.BaseEvent.Creation] when empty
}

// Payload implements eventstore.Command.
func (p *PushedEvent) Payload() any {
	return p
}

func (p *PushedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (p *PushedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	p.BaseEvent = b
}

func NewPushedEvent(
	ctx context.Context,
	aggregate *Aggregate,
	typ Type,
	endpoints []string,
	externalDomain, primaryDomain string,
) *PushedEvent {
	return NewPushedEventWithDate(ctx, aggregate, typ, endpoints, externalDomain, primaryDomain, nil)
}

// NewPushedEventWithDate creates a [PushedEvent] with a fixed Pushed Date.
func NewPushedEventWithDate(
	ctx context.Context,
	aggregate *Aggregate,
	typ Type,
	endpoints []string,
	externalDomain, primaryDomain string,
	pushedDate *time.Time,
) *PushedEvent {
	return &PushedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			&aggregate.Aggregate,
			PushedEventType,
		),
		MilestoneType:  typ,
		Endpoints:      endpoints,
		ExternalDomain: externalDomain,
		PrimaryDomain:  primaryDomain,
		PushedDate:     pushedDate,
	}
}
