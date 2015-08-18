package ga

import (
	"errors"
	"fmt"
	"regexp"
)

var (
	businessID *regexp.Regexp
	designID   *regexp.Regexp
)

func init() {
	businessID = compileRegex("^[A-Za-z0-9\\s\\-_\\.\\(\\)\\!\\?]{1,64}:[A-Za-z0-9\\s\\-_\\.\\(\\)\\!\\?]{1,64}$")
	designID = compileRegex("^[A-Za-z0-9\\s\\-_\\.\\(\\)\\!\\?]{1,64}(:[A-Za-z0-9\\s\\-_\\.\\(\\)\\!\\?]{1,64}){0,4}$")
}

// Event interface TODO
type Event interface {
	Validate() error
}

// DefaultAnnotations defines required information shared by all events
type DefaultAnnotations struct {

	// Use the unique device id if possible. For Android it’s the AID. Should
	// always be the same across game launches.
	UserID string `json:"user_id"`

	// examples: “iPhone6.1”, “GT-I9000”. If not found then “unknown”.
	Device string `json:"device"`

	// Reflects the version of events coming in to the collectors. Current version
	// is 2.
	APIVersion int `json:"v"`

	// Timestamp when the event was created (put in queue/database) on the client.
	// This timestamp should be a corrected one using an offset of time from
	// server_time.
	//
	// 1) The SDK will get the server TS on the init call (each session) and then
	// calculate a difference (within some limit) from the local time and store
	// this ‘offset’.
	//
	// 2) When each event is created it should calculate/adjust the 'client_ts’
	// using the 'offset’.
	ClientTimestamp int `json:"client_ts"`

	// The SDK is submitting events to the servers. For custom solutions ALWAYS
	// use “rest api v2”.
	SDKVersion string `json:"sdk_version"`

	// Operating system version. Like “android 4.4.4”, “ios 8.1”.
	OSVersion string `json:"os_version"`

	// Manufacturer of the hardware the game is played on. Like “apple”,
	// “samsung”, “lenovo”.
	Manufacturer string `json:"manufacturer"`

	// The platform the game is running. Platform is often a subset of os_version
	// like “android”, “windows” etc.
	Platform string `json:"platform"`

	// Universally unique identifier generated on the SDK.
	SessionID string `json:"session_id"`

	// The SDK should count the number of sessions played since it was installed
	// (storing locally and incrementing). The amount should include the session
	// that is about to start.
	SessionNumber uint `json:"session_num"`
}

// NewDefaultAnnotations sets sensible defaults for values on shared annotations.
func NewDefaultAnnotations() *DefaultAnnotations {
	return &DefaultAnnotations{
		Device:       "unknown",
		APIVersion:   APIVersion,
		SDKVersion:   SDKVersion,
		OSVersion:    fmt.Sprintf("webplayer %s", PkgVersion),
		Manufacturer: "unknown",
		Platform:     "webplayer",
	}
}

// User event acts like a session start. It should always be the first event
// in the first batch sent to the collectors and added each time a session
// starts.
type User struct {
	*DefaultAnnotations

	// Category should always be 'user'
	Category string `json:"category"`
}

// NewUserEvent created a new user event with the default annotations
func NewUserEvent(d *DefaultAnnotations) *User {
	return &User{
		Category:           "user",
		DefaultAnnotations: d,
	}
}

// Validate returns an error if any of the event fields is invalid
func (e User) Validate() error {
	if e.Category != "user" {
		return fmt.Errorf("User category MUST be 'user', was %s", e.Category)
	}
	return nil
}

// Business events are for real-money purchases.
type Business struct {
	*DefaultAnnotations

	// Category should always be 'business'
	Category string `json:"category"`

	// A 2 part event id; ItemType:ItemId.
	EventID string `json:"event_id"`

	// The amount of the purchase in cents (integer).
	Amount int `json:"amount"`

	// Currency need to be a 3 letter upper case string to pass validation.
	// In addition the currency need to be a valid currency for correct
	// rate/conversion calculation at a later stage. Look at the following link
	// for a list valid currency values.
	// http://openexchangerates.org/currencies.json.
	Currency string `json:"currency"`

	// Similar to the session_num. Store this value locally and increment each
	// time a business event is submitted during the lifetime (installation) of
	// the game/app.
	TransactionNumber uint `json:"transaction_num"`

	// OPTIONAL
	// A string representing the cart (the location) from which the purchase was
	// made. Could be menu_shop or end_of_level_shop.
	CartType string `json:"cart_type"`
}

// Validate returns an error if any of the event fields is invalid
func (e Business) Validate() error {
	if e.Category != "business" {
		return fmt.Errorf("Business category MUST be 'business', was %s",
			e.Category)
	}
	if ok := businessID.MatchString(e.EventID); !ok {
		return errors.New("EventID doesn't match pattern")
	}
	if e.Currency == "" {
		return errors.New("Currency missing")
	}
	found := false
	for _, c := range Currencies {
		if c == e.Currency {
			found = true
			break
		}
	}
	if !found {
		return errors.New("Currency is invalid. Check Currencies for a valid list")
	}
	return nil
}

// NewBusinessEvent created a new user event with the default annotations
func NewBusinessEvent(d *DefaultAnnotations) *Business {
	return &Business{
		Category:           "business",
		DefaultAnnotations: d,
	}
}

// SessionEnd event should always be sent whenever a session is determined
// to be over. For example whenever a mobile device is ‘going-to-background’ or
// when a user quit your game in other ways. Only one SessionEnd event per
// session should be generated/sent.
type SessionEnd struct {
	*DefaultAnnotations

	// Category should always be 'session_end'
	Category string `json:"category"`

	// Session length in seconds.
	Length int `json:"length"`
}

// NewSessionEndEvent created a new user event with the default annotations
func NewSessionEndEvent(d *DefaultAnnotations) *SessionEnd {
	return &SessionEnd{
		Category:           "session_end",
		DefaultAnnotations: d,
	}
}

// Validate returns an error for invalid length
func (e SessionEnd) Validate() error {
	if e.Category != "session_end" {
		return fmt.Errorf("SessionEnd category MUST be 'session_end', was %s",
			e.Category)
	}
	if e.Length < 0 {
		return fmt.Errorf("Length must be equal or greater than 0 (%d)", e.Length)
	}
	return nil
}

/*
// Design events are for general in-game events that are not covered by other
// events.
type Design struct {
	*DefaultAnnotations

	// A 1-5 part event id.
	EventID string `json:"event_id"`

	// Optional value. float.
	Value float64 `json:"value,omitempty"`
}

// Validate returns an error if any of the event fields is invalid
func (e Design) Validate() error {
	if ok := designID.MatchString(e.EventID); !ok {
		return errors.New("EventID doesn't match pattern")
	}
	return nil
}
*/

func compileRegex(pattern string) *regexp.Regexp {
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}
	return re
}
