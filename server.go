package ga

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"time"
)

const (
	// PkgVersion is the current version of this package. Follows major, minor and
	// patch conventions
	PkgVersion = "0.0.2"

	// APIVersion is the current version supported by GameAnalytics
	APIVersion = 2

	// SDKVersion is the current version supported by GameAnalytics
	SDKVersion = "rest api v2"

	// InitRoute is the url part for the init request
	InitRoute = "init"

	// EventsRoute is the url part for events request
	EventsRoute = "events"

	// SandboxGameKey is the game key for the GameAnalytics sandbox-api
	SandboxGameKey = "5c6bcb5402204249437fb5a7a80a4959"

	// SandboxSecretKey is the secret key for the GameAnalytics sandbox-api
	SandboxSecretKey = "16813a12f718bc5c620f56944e1abc3ea13ccbac"
)

// APIStatus is the GameAnalytics response of the init event. If Enabled is
// false, the server shouldn't send any events.
type APIStatus struct {
	Enabled         bool
	ServerTimestamp int `json:"server_ts"`
	Flags           []string
}

// Server wraps the API endpoint and allows events to be sent
type Server struct {
	// GameKey provided by GameAnalytics for the account
	GameKey string `json:"-"`

	// SecretKey provided by GameAnalytics for the account
	SecretKey string `json:"-"`

	// URL endpoint for GameAnalytics API
	URL string `json:"-"`

	// Platform represents the platform of the SDK
	Platform string `json:"platform"`

	// OSVersion represents the Operational System Version of the SDK
	OSVersion string `json:"os_version"`

	// SDKVersion is the version of the SDK
	SDKVersion string `json:"sdk_version"`

	// Offset from GameAnalytics API and this server
	TimestampOffset int `json:"-"`

	APIStatus
}

// NewServer returns a server with default values for the GameAnalytics
// custom SDK implementation.
func NewServer(gameKey, secretKey string) *Server {
	return &Server{
		URL:        fmt.Sprintf("http://api.gameanalytics.com/v%d", APIVersion),
		SDKVersion: SDKVersion,
		OSVersion:  runtime.Version(),
		Platform:   "go",
		GameKey:    gameKey,
		SecretKey:  secretKey,
	}
}

// NewSandboxServer return a server with default values for the GameAnalytics
// sandbox API
func NewSandboxServer() *Server {
	return &Server{
		URL:        fmt.Sprintf("http://sandbox-api.gameanalytics.com/v%d", APIVersion),
		SDKVersion: SDKVersion,
		OSVersion:  runtime.Version(),
		Platform:   "go",
		GameKey:    SandboxGameKey,
		SecretKey:  SandboxSecretKey,
	}
}

// Start does the initial request to GameAnalytics API
func (s *Server) Start() error {
	payload, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("Init marshal payload failed (%v)", err)
	}

	body, err := s.post(InitRoute, payload)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(body), &s.APIStatus)
	if err != nil {
		return fmt.Errorf("APIStatus unmarshal failed (%v)", err)
	}

	if !s.Enabled {
		return fmt.Errorf("API is disabled. Server can't send any events")
	}

	epoch := int(time.Now().Unix())
	s.TimestampOffset = s.ServerTimestamp - epoch
	return nil
}

// SendEvent posts a single event to GameAnalytics using the server config
func (s *Server) SendEvent(e Event) error {
	return s.SendEvents([]Event{e})
}

// SendEvents posts one or more events to GameAnalytics using the server config
func (s *Server) SendEvents(e []Event) error {
	payload, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("Init marshal payload failed (%v)", err)
	}

	result, err := s.post(EventsRoute, payload)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Event sent (%s), response: %s\n", payload, result)

	return nil
}

// Post sends a payload using the server config
func (s *Server) post(route string, payload []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", s.URL, s.GameKey, route)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("Preparing request failed (%v)", err)
	}

	auth := computeHmac256(payload, s.SecretKey)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "application/json") //TODO add gzip compression

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Server request failed (%v)", err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Expected status code 200, got %d. Body: %s",
			res.StatusCode, body)
	}
	return []byte(body), nil
}

// computeHmac256 returns the raw body content from the request using the secret
// key (private key) as the hashing key and then encoding it using base64.
func computeHmac256(payload []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(payload))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
