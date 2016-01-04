package main

import (
	"fmt"
	"os"

	"github.com/iu0v1/daslog"
)

func main() {
	o := daslog.Options{
		Destination: os.Stdout,
		Prefix:      "{{.O}} [{{.Q}}]: ",
		LogLevel:    daslog.UrgencyLevelCritical,
	}

	l, err := daslog.New(o)
	if err != nil {
		fmt.Print(err)
		return
	}

	// notice in Log style
	l.Log(daslog.UrgencyLevelNotice, "test notice message")

	// info
	l.Info("test info message")

	// error
	l.Errorf("%s %s %s", "test", "error", "message")
}
