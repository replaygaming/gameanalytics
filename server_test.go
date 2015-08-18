package ga

import (
	"reflect"
	"runtime"
	"testing"
)

func TestNewServer(t *testing.T) {
	result := NewServer("key", "secret")
	expected := &Server{
		URL:        "http://api.gameanalytics.com/v2",
		SDKVersion: "rest api v2",
		OSVersion:  runtime.Version(),
		Platform:   "go",
		GameKey:    "key",
		SecretKey:  "secret",
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expecting %v to be equal to %v", result, expected)
	}
}

func TestNewSandboxServer(t *testing.T) {
	result := NewSandboxServer()
	expected := &Server{
		URL:        "http://sandbox-api.gameanalytics.com/v2",
		SDKVersion: "rest api v2",
		OSVersion:  runtime.Version(),
		Platform:   "go",
		GameKey:    SandboxGameKey,
		SecretKey:  SandboxSecretKey,
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expecting %v to be equal to %v", result, expected)
	}
}

func TestComputeHmac256(t *testing.T) {
	result := computeHmac256([]byte("test"), "secret")
	expected := "Aymga2LNFrM+tnkr6MYLFY2Jou46h2/Omogeu0iMCRQ="
	if result != expected {
		t.Errorf("Expecting %s to be equal to %s", result, expected)
	}
}
