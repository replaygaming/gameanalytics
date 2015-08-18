package ga

import "testing"

type eventCase struct {
	event Event
	valid bool
}

func TestNewUserEvent(t *testing.T) {
	s := &DefaultAnnotations{}
	e := NewUserEvent(s)
	if e.Category != "user" {
		t.Errorf("Expected user event to has user category (%s)", e.Category)
	}
}

func validateEvents(cases []eventCase, t *testing.T) {
	for _, c := range cases {
		err := c.event.Validate()
		if c.valid {
			if err != nil {
				t.Errorf("Expected event (%v) to be valid", c.event)
			}
		} else {
			if err == nil {
				t.Errorf("Expected event (%v) to be invalid", c.event)
			}
		}
	}
}

func TestUser_Validate(t *testing.T) {
	var cases = []eventCase{
		{&User{Category: "other"}, false},
		{&User{Category: "user"}, true},
	}
	validateEvents(cases, t)
}

func TestNewSessionEndEvent(t *testing.T) {
	s := &DefaultAnnotations{}
	e := NewSessionEndEvent(s)
	if e.Category != "session_end" {
		t.Errorf("Expected session_end event to has 'session_end' category (%s)",
			e.Category)
	}
}

func TestSessionEnd_Validate(t *testing.T) {
	var cases = []eventCase{
		{&SessionEnd{Category: "other"}, false},
		{&SessionEnd{Category: "session_end", Length: -1}, false},
		{&SessionEnd{Category: "session_end", Length: 0}, true},
	}
	validateEvents(cases, t)
}

func TestNewBusinessEvent(t *testing.T) {
	s := &DefaultAnnotations{}
	e := NewBusinessEvent(s)
	if e.Category != "business" {
		t.Errorf("Expected business event to has business category (%s)", e.Category)
	}
}

func TestBusiness_Validate(t *testing.T) {
	var cases = []eventCase{
		{&Business{Category: "other"}, false},
		{&Business{Category: "business", EventID: "WrongPattern"}, false},
		{&Business{Category: "business", EventID: "Correct:Pattern", Currency: ""}, false},
		{&Business{Category: "business", EventID: "Correct:Pattern", Currency: "AAA"}, false},
		{&Business{Category: "business", EventID: "Correct:Pattern", Currency: "USD"}, true},
	}
	validateEvents(cases, t)
}
