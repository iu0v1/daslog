package daslog

import (
	"bytes"
	"testing"
)

// almost dummy test
// TODO : add tests
func TestMain(t *testing.T) {
	var buf bytes.Buffer

	o := Options{
		Destination: &buf,
		Prefix:      "{{.O}} [{{.Q}}]: ",
		LogLevel:    UrgencyLevelCritical,
	}

	l, err := New(o)
	if err != nil {
		t.Fatalf("%v", err)
	}

	// notice in Log style
	l.Log(UrgencyLevelNotice, "test notice message")

	// info
	l.Info("test info message")

	// error
	l.Errorf("%s %s %s", "test", "error", "message")
}
