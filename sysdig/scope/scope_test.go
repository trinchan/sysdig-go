package scope

import (
	"testing"
)

func TestScopeString(t *testing.T) {
	tests := []struct {
		name  string
		scope *Scope
		want  string
	}{
		{
			name:  "Nil",
			scope: nil,
			want:  "",
		},
		{
			name:  "Empty",
			scope: New(),
			want:  "",
		},
		{
			name:  "Is",
			scope: New().AddIsSelection("foo", "bar"),
			want:  `foo = 'bar'`,
		},
		{
			name:  "IsNot",
			scope: New().AddIsNotSelection("foo", "bar"),
			want:  `foo != 'bar'`,
		},
		{
			name:  "Contains",
			scope: New().AddContainsSelection("foo", "bar"),
			want:  `foo contains 'bar'`,
		},
		{
			name:  "NotIn",
			scope: New().AddDoesNotContainSelection("foo", "bar"),
			want:  `not foo contains 'bar'`,
		},
		{
			name:  "StartsWith",
			scope: New().AddStartsWithSelection("foo", "bar"),
			want:  `foo starts with 'bar'`,
		},
		{
			name:  "In",
			scope: New().AddInSelection("foo", "bar"),
			want:  `foo in ('bar')`,
		},
		{
			name:  "NotIn",
			scope: New().AddNotInSelection("foo", "bar"),
			want:  `not foo in ('bar')`,
		},
		{
			name:  "InMultiple",
			scope: New().AddInSelection("foo", "bar", "baz"),
			want:  `foo in ('bar', 'baz')`,
		},
		{
			name:  "NotInMultiple",
			scope: New().AddNotInSelection("foo", "bar", "baz"),
			want:  `not foo in ('bar', 'baz')`,
		},
		{
			name:  "MultipleIs",
			scope: New().AddIsSelection("foo", "bar").AddIsSelection("baz", "biz"),
			want:  `foo = 'bar' and baz = 'biz'`,
		},
		{
			name:  "MultipleIsAndContains",
			scope: New().AddIsSelection("a", "b").AddIsSelection("c", "d").AddInSelection("e", "f"),
			want:  `a = 'b' and c = 'd' and e in ('f')`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.scope.String()
			if got != test.want {
				t.Errorf("got: %s, want: %s", got, test.want)
			}
		})
	}
}

func TestEventScopeString(t *testing.T) {
	tests := []struct {
		name  string
		scope *EventScope
		want  string
	}{
		{
			name:  "Empty",
			scope: NewEventScope(),
			want:  "",
		},
		{
			name:  "WithLabels",
			scope: NewEventScopeWithLabels(map[string]string{"foo": "bar", "baz": "biz"}),
			want:  `foo = 'bar' and baz = 'biz'`,
		},
		{
			name:  "Is",
			scope: NewEventScope().AddIsSelection("foo", "bar"),
			want:  `foo = 'bar'`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.scope.String()
			if got != test.want {
				t.Errorf("got: %s, want: %s", got, test.want)
			}
		})
	}
}
