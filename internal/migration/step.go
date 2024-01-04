package migration

import "github.com/zitadel/zitadel/internal/eventstore"

var _ eventstore.QueryReducer = (*StepStates)(nil)

type Step struct {
	*SetupStep

	state StepState
}

type StepStates struct {
	eventstore.ReadModel
	Steps []*Step
}

// Query implements eventstore.QueryReducer.
func (*StepStates) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(aggregateType).
		AggregateIDs(aggregateID).
		EventTypes(StartedType, doneType, repeatableDoneType, failedType).
		Builder()
}

// Reduce implements eventstore.QueryReducer.
func (s *StepStates) Reduce() error {
	for _, event := range s.Events {
		step := event.(*SetupStep)
		state := s.byName(step.Name)
		if state == nil {
			state = new(Step)
			s.Steps = append(s.Steps, state)
		}
		state.SetupStep = step
		switch step.EventType {
		case StartedType:
			state.state = StepStarted
		case doneType:
			state.state = StepDone
		case repeatableDoneType:
			state.state = StepDone
		case failedType:
			state.state = StepFailed
		}
	}
	return s.ReadModel.Reduce()
}

func (s *StepStates) byName(name string) *Step {
	for _, step := range s.Steps {
		if step.Name == name {
			return step
		}
	}
	return nil
}

func (s *StepStates) lastByState(stepState StepState) (step *Step) {
	for _, state := range s.Steps {
		if state.state != stepState {
			continue
		}
		if step == nil {
			step = state
			continue
		}
		if step.CreatedAt().After(state.CreatedAt()) {
			continue
		}

		step = state
	}
	return step
}

func (s *StepStates) byState(name string) (step *Step) {
	for _, state := range s.Steps {
		if state.Name != name {
			continue
		}
		return state
	}
	return nil
}

type StepState int32

const (
	StepStarted StepState = iota
	StepDone
	StepFailed
)
