# DASLOG

Daslog is a simple logger service for Go (golang) programms.

**DA**mn **S**imple **LOG**ger.

```
go get -u github.com/iu0v1/daslog
```

### Why?
Yes. This is yet another logger. I made this package, because any other logger that I met - is a monster. They can do a lot of tricks, with one exception - a simple and transparent logging. Sometimes, you need only logging, without the possibility of launching a space shuttle through the services of Amazon or Google.

### Prefix format
You can use some `GNU date` format commands in prefix line.
```
'GNU date' compatible format commands:

{{.F}} - full date                               (2016-01-04)
{{.T}} - time                                    (16:52:36)
{{.r}} - 12-hour clock time                      (11:11:04 PM)
{{.Y}} - year                                    (2016)
{{.y}} - last two digits of year                 (00..99)
{{.m}} - month                                   (01..12)
{{.b}} - locale's abbreviated month name         (Jan)
{{.B}} - locale's full month name                (January)
{{.d}} - day of month                            (01)
{{.a}} - locale's abbreviated weekday name       (Sun)
{{.A}} - locale's full weekday name              (Sunday)
{{.H}} - 24-hour clock hour                      (00..23)
{{.I}} - 12-hour clock hour                      (01..12)
{{.M}} - minute                                  (00..59)
{{.S}} - second                                  (00..60)
{{.p}} - locale's equivalent of either AM or PM  (AM)

Incompatible format commands:

{{.O}} - same as {{.F}} + {{.T}}                 (2016-01-04 16:52:36)
{{.Q}} - message urgency level                   (info, critical, etc)

```

### Simple example
```
package main

import (
	"fmt"
	"os"

	"github.com/iu0v1/daslog"
)

func main() {
	o := Options{
		Destination: os.Stdout,
		Prefix:      "{{.O}} [{{.Q}}]: ",
		LogLevel:    UrgencyLevelCritical,
	}

	l, err := New(o)
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
```
Output:
```
2016-01-04 21:16:41 [notice]: test notice message
2016-01-04 21:16:41 [info]: test info message
2016-01-04 21:16:42 [error]: test error message
```

### DOC
For more infomation, please look at the [examples](https://github.com/iu0v1/daslog/tree/master/example) and read the [doc](http://godoc.org/github.com/iu0v1/daslog).