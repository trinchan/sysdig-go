package scope

import (
	"fmt"
	"strings"
)

// Scope defines a filter for an EventsService.ListEvents.
type Scope struct {
	selections []scopeSelection
}

// EventScope defines the Scope labels for an EventsService.CreateEvent
type EventScope Scope

// Selector defines a type for Scope filter operators.
type Selector string

const (
	// SelectionIs filters EventsService.ListEvents to Events which do not exactly match the provided value.
	// It is also used to set Scope labels for a created Event.
	SelectionIs Selector = "="
	// SelectionIsNot filters EventsService.ListEvents to Events which do not exactly match the provided value.
	SelectionIsNot Selector = "!="
	// SelectionIn filters EventsService.ListEvents to Events which exactly match one of the provided values.
	SelectionIn Selector = "in"
	// SelectionNotIn filters EventsService.ListEvents to Events which do not exactly match one of the provided values.
	SelectionNotIn Selector = "not in"
	// SelectionContains filters EventsService.ListEvents to Events which substring match the provided value.
	SelectionContains Selector = "contains"
	// SelectionDoesNotContain filters EventsService.ListEvents to Events which do not substring match the provided value.
	SelectionDoesNotContain Selector = "does not contain"
	// SelectionStartsWith filters EventsService.ListEvents to Events which are prefixed with the provided value.
	SelectionStartsWith Selector = "starts with"
)

type scopeSelection struct {
	label    string
	values   []string
	selector Selector
}

// New initializes a new Scope to be used for filters in EventsService.ListEvents.
func New() *Scope {
	return &Scope{}
}

// NewEventScope initializes a new EventScope to be used for adding scope labels to an event in EventsService.CreateEvent.
func NewEventScope() *EventScope {
	return NewEventScopeWithLabels(nil)
}

// NewEventScopeWithLabels initializes a new EventScope with the provided scope labels.
func NewEventScopeWithLabels(labels map[string]string) *EventScope {
	s := &EventScope{}
	for k, v := range labels {
		s.AddIsSelection(k, v)
	}
	return s
}

// AddIsSelection adds a label to the EventScope with the given value.
func (s *EventScope) AddIsSelection(label, value string) *EventScope {
	(*Scope)(s).AddIsSelection(label, value)
	return s
}

// String is implemented via the Scope.String() method.
func (s *EventScope) String() string {
	return (*Scope)(s).String()
}

// AddSelection adds a filter to the Scope with the given label, value and selector.
// Note: You should prefer to use the explicitly scoped Add functions to ensure you do not
// accidentally add multiple values when only one is supported. e.g. if you add multiple
// values to SelectionIs Selection, the values will be joined which is probably
// not what you want.
func (s *Scope) AddSelection(selector Selector, label, value string) *Scope {
	s.selections = append(s.selections, scopeSelection{
		label:    label,
		values:   []string{value},
		selector: selector,
	})
	return s
}

// AddSelectionMultiple adds a filter to the Scope with the given label, value and selector.
// // Note: You should prefer to use the explicitly scoped Add functions to ensure you do not
// // accidentally add multiple values when only one is supported. e.g. if you add multiple
// // values to SelectionIs Selection, the values will be joined which is probably
// // not what you want.
func (s *Scope) AddSelectionMultiple(selector Selector, label string, values ...string) *Scope {
	s.selections = append(s.selections, scopeSelection{
		label:    label,
		values:   values,
		selector: selector,
	})
	return s
}

// AddContainsSelection adds a filter to the Scope with the given label, value and SelectionContains selector.
func (s *Scope) AddContainsSelection(label string, value string) *Scope {
	return s.AddSelectionMultiple(SelectionContains, label, value)
}

// AddDoesNotContainSelection adds a filter to the Scope with the given label, value and DoesNotContain selector.
func (s *Scope) AddDoesNotContainSelection(label string, value string) *Scope {
	return s.AddSelectionMultiple(SelectionDoesNotContain, label, value)
}

// AddIsSelection adds a filter to the Scope with the given label, value and SelectionIs selector.
func (s *Scope) AddIsSelection(label, value string) *Scope {
	return s.AddSelection(SelectionIs, label, value)
}

// AddIsNotSelection adds a filter to the Scope with the given label, value and SelectionIsNot selector.
func (s *Scope) AddIsNotSelection(label, value string) *Scope {
	return s.AddSelection(SelectionIsNot, label, value)
}

// AddInSelection adds a filter to the Scope with the given label, value and SelectionIn selector.
func (s *Scope) AddInSelection(label string, values ...string) *Scope {
	return s.AddSelectionMultiple(SelectionIn, label, values...)
}

// AddNotInSelection adds a filter to the Scope with the given label, value and SelectionNotIn selector.
func (s *Scope) AddNotInSelection(label string, values ...string) *Scope {
	return s.AddSelectionMultiple(SelectionNotIn, label, values...)
}

// AddStartsWithSelection adds a filter to the Scope with the given label, value and SelectionStartsWith selector.
func (s *Scope) AddStartsWithSelection(label, value string) *Scope {
	return s.AddSelection(SelectionStartsWith, label, value)
}

// String defines fmt.Stringer for Scope. It converts it to the Sysdig format for Scope strings.
func (s *Scope) String() string {
	if s == nil {
		return ""
	}
	var b strings.Builder
	for i, selection := range s.selections {
		if i > 0 {
			b.WriteString(" and ")
		}
		var prefix string
		selector := selection.selector
		// First set the "not" prefix in front of the Scope filter and reverse the selection to parse correctly.
		switch selection.selector {
		case SelectionDoesNotContain:
			prefix = "not "
			selector = SelectionContains
		case SelectionNotIn:
			prefix = "not "
			selector = SelectionIn
		}

		// If someone misuses the client and sets multiple values for a Selection that isn't SelectionIn or SelectionNotIn,
		// then it'll join them together and still work. Would probably be better to validate and fail though.
		var joined strings.Builder
		for c, v := range selection.values {
			if c > 0 {
				joined.WriteString(", ")
			}
			joined.WriteRune('\'')
			joined.WriteString(v)
			joined.WriteRune('\'')
		}
		switch selection.selector {
		case SelectionIn, SelectionNotIn:
			b.WriteString(fmt.Sprintf(`%s%s %s (%s)`, prefix, selection.label, selector, joined.String()))
		default:
			b.WriteString(fmt.Sprintf(`%s%s %s %s`, prefix, selection.label, selector, joined.String()))
		}
	}
	return b.String()
}
