package ga

import "testing"

func TestNewUserEvent(t *testing.T) {
	s := &DefaultAnnotations{}
	e := NewUserEvent(s)
	if e.Category != "user" {
		t.Errorf("Expected user event to has user category (%s)", e.Category)
	}
}

func TestUser_Validate(t *testing.T) {
	e := &User{Category: "other"}
	err := e.Validate()
	if err == nil {
		t.Errorf("Expected user category to be invalid (%s)", e.Category)
	}

	e.Category = "user"
	err = e.Validate()
	if err != nil {
		t.Errorf("Expected user category to be valid (%v)", e)
	}
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
	e := &SessionEnd{Category: "other"}
	err := e.Validate()
	if err == nil {
		t.Errorf("Expected session_end category to be invalid (%s)", e.Category)
	}
	e.Category = "session_end"

	e.Length = -1
	err = e.Validate()
	if err == nil {
		t.Errorf("Expected session_end length to be invalid (%d)", e.Length)
	}

	e.Length = 0
	err = e.Validate()
	if err != nil {
		t.Errorf("Expected session_end category to be valid (%v)", e)
	}
}

func TestNewBusinessEvent(t *testing.T) {
	s := &DefaultAnnotations{}
	e := NewBusinessEvent(s)
	if e.Category != "business" {
		t.Errorf("Expected business event to has business category (%s)", e.Category)
	}
}

func TestBusiness_Validate(t *testing.T) {
	e := &Business{Category: "other"}
	err := e.Validate()
	if err == nil {
		t.Errorf("Expected business category to be invalid (%s)", e.Category)
	}
	e.Category = "business"

	e.EventID = "WrongPattern"
	err = e.Validate()
	if err == nil {
		t.Errorf("Expected business event_id to be invalid (%s)", e.EventID)
	}
	e.EventID = "Correct:Pattern"

	e.Currency = "AAA"
	err = e.Validate()
	if err == nil {
		t.Errorf("Expected business currency to be invalid (%s)", e.Currency)
	}
	e.Currency = "USD"

	err = e.Validate()
	if err != nil {
		t.Errorf("Expected business to be valid (%v)", e)
	}
}
