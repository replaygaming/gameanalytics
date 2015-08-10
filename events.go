package ga

import (
	"errors"
	"fmt"
	"log"
	"regexp"
)

var (
	businessID *regexp.Regexp
	designID   *regexp.Regexp
)

func init() {
	var err error
	businessID, err = regexp.Compile(
		"^[A-Za-z0-9\\s\\-_\\.\\(\\)\\!\\?]{1,64}:[A-Za-z0-9\\s\\-_\\.\\(\\)\\!\\?]{1,64}$")
	if err != nil {
		log.Fatal(err)
	}

	designID, err = regexp.Compile(
		"^[A-Za-z0-9\\s\\-_\\.\\(\\)\\!\\?]{1,64}(:[A-Za-z0-9\\s\\-_\\.\\(\\)\\!\\?]{1,64}){0,4}$")
	if err != nil {
		log.Fatal(err)
	}
}

// Event interface TODO
type Event interface{}

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

// DefaultsRequiredValues sets sensible defaults for missing values on shared
// annotations.
func (a *DefaultAnnotations) DefaultsRequiredValues() {
	if a.Device == "" {
		a.Device = "unknown"
	}
	if a.APIVersion == 0 {
		a.APIVersion = APIVersion
	}
	if a.SDKVersion == "" {
		a.SDKVersion = SDKVersion
	}
	if a.OSVersion == "" {
		a.OSVersion = fmt.Sprintf("webplayer %s", PkgVersion)
	}
	if a.Manufacturer == "" {
		a.Manufacturer = "unknown"
	}
	if a.Platform == "" {
		a.Platform = "webplayer"
	}
}

// User event acts like a session start. It should always be the first event
// in the first batch sent to the collectors and added each time a session
// starts.
type User struct {
	*DefaultAnnotations

	// Category is always 'user'
	Category string `json:"category"`
}

// NewUserEvent created a new user event with the default annotations
func NewUserEvent(d *DefaultAnnotations) *User {
	return &User{
		Category:           "user",
		DefaultAnnotations: d,
	}
}

// ValidateAttributes always returns nil for user
func (User) ValidateAttributes() error {
	return nil
}

// Business events are for real-money purchases.
type Business struct {
	*DefaultAnnotations

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
}

// ValidateAttributes returns whether the event fields are valid
func (e Business) ValidateAttributes() error {
	if ok := businessID.MatchString(e.EventID); !ok {
		return errors.New("EventID doesn't match pattern")
	}
	//TODO validate currency
	return nil
}

// SessionEnd event should always be sent whenever a session is determined
// to be over. For example whenever a mobile device is ‘going-to-background’ or
// when a user quit your game in other ways. Only one SessionEnd event per
// session should be generated/sent.
type SessionEnd struct {
	*DefaultAnnotations

	// Category is always 'session_end'
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

// ValidateAttributes returns an error for invalid length
func (e SessionEnd) ValidateAttributes() error {
	if e.Length < 0 {
		return errors.New("Length must be equal or greater than 0")
	}
	return nil
}

// Design events are for general in-game events that are not covered by other
// events.
type Design struct {
	*DefaultAnnotations

	// A 1-5 part event id.
	EventID string `json:"event_id"`

	// Optional value. float.
	Value float64 `json:"value,omitempty"`
}

// ValidateAttributes returns whether the event fields are valid
func (e Design) ValidateAttributes() error {
	if ok := designID.MatchString(e.EventID); !ok {
		return errors.New("EventID doesn't match pattern")
	}
	return nil
}
