package daslog

import (
	"bytes"
	"fmt"
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
		fmt.Print(err)
		return
	}

	// notice in Log style
	l.Log(UrgencyLevelNotice, "test notice message")

	// info
	l.Info("test info message")

	// error
	l.Errorf("%s %s %s", "test", "error", "message")
}
