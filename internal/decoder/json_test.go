package decoder

import (
	"encoding/json"
	"testing"
)

func TestJson_Single(t *testing.T) {
	testData := []byte(`{"test": 1234, "inner": {"test": 1234}}`)

	var m Messages

	if err := json.Unmarshal(testData, &m); err != nil {
		t.Error(err)
	}

	if len(m) != 1 {
		t.Errorf("expected 1 message got %d", len(m))
	}
}

func TestJson_Multi(t *testing.T) {
	testData := []byte(`[{"test": 1234, "inner": {"test": 1234}}, {"test2": "XXX"}]`)

	var m Messages

	if err := json.Unmarshal(testData, &m); err != nil {
		t.Error(err)
	}

	if len(m) != 2 {
		t.Errorf("expected 2 messages got %d", len(m))
	}
}
