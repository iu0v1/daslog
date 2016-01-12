package daslog

import (
	"bytes"
	"log"
	"testing"
)

func BenchmarkLog(b *testing.B) {
	var buf bytes.Buffer
	l := log.New(&buf, "test", log.LstdFlags)

	for i := 0; i < b.N; i++ {
		l.Printf("test %s", "message")
	}
}

func BenchmarkDaslogPurePrefix(b *testing.B) {
	var buf bytes.Buffer
	o := Options{
		Destination: &buf,
		Prefix:      "test",
		LogLevel:    UrgencyLevelCritical,
	}

	l, _ := New(o)

	for i := 0; i < b.N; i++ {
		l.Infof("test %s", "message")
	}
}

// template prefix without a {{.Q}} option
func BenchmarkDaslogTemplatePrefix(b *testing.B) {
	var buf bytes.Buffer
	o := Options{
		Destination: &buf,
		Prefix:      "{{.O}}: ",
		LogLevel:    UrgencyLevelCritical,
	}

	l, _ := New(o)

	for i := 0; i < b.N; i++ {
		l.Infof("test %s", "message")
	}
}

// template prefix with a {{.Q}} option
func BenchmarkDaslogTemplatePrefixQ(b *testing.B) {
	var buf bytes.Buffer
	o := Options{
		Destination: &buf,
		Prefix:      "{{.O}} [{{.Q}}]: ",
		LogLevel:    UrgencyLevelCritical,
	}

	l, _ := New(o)

	for i := 0; i < b.N; i++ {
		l.Infof("test %s", "message")
	}
}
