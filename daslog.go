// Package daslog is a simple logger service.
//
// Explanation of the name: DAmn Simple LOGger.
//
// Simple example:
//
//      func main() {
//      	o := daslog.Options{
//              Destination: os.Stdout,
//              Prefix:      "{{.O}} [{{.Q}}]: ",
//              LogLevel:    daslog.UrgencyLevelCritical,
//      	}
//
//      	l, err := daslog.New(o)
//      	if err != nil {
//      		fmt.Print(err)
//      		return
//      	}
//
//      	l.Log(daslog.UrgencyLevelNotice, "test notice message")
//
//      	l.Info("test info message")
//
//      	l.Errorf("%s %s %s", "test", "error", "message")
//      }
//
// Output:
//      2016-01-04 21:16:41 [notice]: test notice message
//      2016-01-04 21:16:41 [info]: test info message
//      2016-01-04 21:16:42 [error]: test error message
//
package daslog

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"
)

// UrgencyLevel - declare the level of informativity of log message.
type UrgencyLevel int

// LogHandler - log handler type.
type LogHandler func(urgencyLevel UrgencyLevel, message string)

// predefined UrgencyLevel levels
const (
	UrgencyLevelNone UrgencyLevel = iota
	UrgencyLevelNotice
	UrgencyLevelInfo
	UrgencyLevelError
	UrgencyLevelCritical
)

// Options - struct, which is used to configure a Daslog.
type Options struct {
	// LogLevel provides the opportunity to choose the level of
	// information messages.
	// Each level includes the messages from the previous level.
	// UrgencyLevelNone     - no messages // 0
	// UrgencyLevelNotice   - notice      // 1
	// UrgencyLevelInfo     - info        // 2
	// UrgencyLevelError    - error       // 3
	// UrgencyLevelCritical - critical    // 4
	//
	// Default: 'UrgencyLevelNone'.
	LogLevel UrgencyLevel

	// Prefix provides the ability to add prefix text to all log messages.
	// In addition, you can use commands to display the date and time (these
	// format commands are based on the linux 'GNU date' format syntax).
	//
	// 'GNU date' compatible format commands:
	//
	// {{.F}} - full date                               (2016-01-04)
	// {{.T}} - time                                    (16:52:36)
	// {{.r}} - 12-hour clock time                      (11:11:04 PM)
	//
	// {{.Y}} - year                                    (2016)
	// {{.y}} - last two digits of year                 (00..99)
	// {{.m}} - month                                   (01..12)
	// {{.b}} - locale's abbreviated month name         (Jan)
	// {{.B}} - locale's full month name                (January)
	// {{.d}} - day of month                            (01)
	// {{.a}} - locale's abbreviated weekday name       (Sun)
	// {{.A}} - locale's full weekday name              (Sunday)
	// {{.H}} - 24-hour clock hour                      (00..23)
	// {{.I}} - 12-hour clock hour                      (01..12)
	// {{.M}} - minute                                  (00..59)
	// {{.S}} - second                                  (00..60)
	// {{.p}} - locale's equivalent of either AM or PM  (AM)
	//
	// Incompatible format commands:
	//
	// {{.O}} - same as {{.F}} + {{.T}}                 (2016-01-04 16:52:36)
	// {{.Q}} - message urgency level                   (info, critical, etc)
	Prefix string

	// Destination provides the opportunity to choose the own
	// destination for log messages (errors, info, etc).
	//
	// Default: 'os.Stdout'.
	Destination io.Writer

	// LogHandler takes a log messages to bypass the internal
	// mechanism of the message processing.
	//
	// If LogHandler is selected - all log settings will be ignored.
	Handler LogHandler
}

type prefixTemplateData struct {
	F, T, Or, Y, Oy, Om string
	Ob, B, Od, Oa, A, H string
	I, M, S, Op, O, Q   string
}

// Daslog - main struct.
type Daslog struct {
	options *Options

	templateInPrefix  bool
	timeTemplate      string
	urgencyInTemplate bool

	urgencyList []string
}

// New - return a new copy of Daslog.
func New(options Options) (*Daslog, error) {
	l := &Daslog{
		options: &options,
	}

	if l.options.Handler == nil {
		l.options.Handler = l.logHandler
	}

	if l.options.Destination == nil {
		l.options.Destination = os.Stdout
	}

	if l.options.Prefix != "" {
		if t, _ := regexp.MatchString("{{\\s*\\.\\w\\s*}}", l.options.Prefix); t {
			// repleace all unexported letters and exceptions
			replaceList := []string{"r", "y", "m", "b", "d", "a", "p", "Q"}
			newPrefix := l.options.Prefix

			for _, r := range replaceList {
				a := regexp.MustCompile("{{\\s*\\." + r + "\\s*}}")
				if r == "Q" {
					newPrefix = a.ReplaceAllString(newPrefix, "<{Q}>")
					l.urgencyInTemplate = true
					continue
				}
				newPrefix = a.ReplaceAllString(newPrefix, "{{.O"+r+"}}")
			}

			l.options.Prefix = newPrefix

			// prepare template
			tmpl, err := template.New("prefix").Parse(l.options.Prefix)
			if err != nil {
				e := regexp.MustCompile("template: prefix:1")
				return l, fmt.Errorf("%s\n", e.ReplaceAllString(err.Error(), "daslog: prefix error"))
			}
			l.templateInPrefix = true

			// check unknown format variables and errors
			var buf bytes.Buffer
			dummyData := prefixTemplateData{}
			if err := tmpl.Execute(&buf, &dummyData); err != nil {
				if e, _ := regexp.MatchString("is not a field of", err.Error()); e {
					e := regexp.MustCompile("<\\..*>").FindAllString(err.Error(), -1)
					ev := regexp.MustCompile("(<|>)").ReplaceAllString(e[0], "")
					return l, fmt.Errorf("daslog: unknown format variable in prefix: {{%s}}\n", ev)
				}
				return l, err
			}

			// call template
			buf.Reset()

			data := prefixTemplateData{
				F:  "2006-01-02",
				T:  "15:04:05",
				Or: "3:04:05 PM",
				Y:  "2006",
				Oy: "06",
				Om: "01",
				Ob: "Jan",
				B:  "January",
				Od: "02",
				Oa: "Mon",
				A:  "Monday",
				H:  "15",
				I:  "03",
				M:  "04",
				S:  "05",
				Op: "PM",

				O: "2006-01-02 15:04:05",
			}

			tmpl.Execute(&buf, &data)

			l.timeTemplate = buf.String()

			l.urgencyList = []string{
				"none",
				"notice",
				"info",
				"error",
				"critical",
			}
		}
	}

	return l, nil
}

// logHandler - default message handler.
func (l *Daslog) logHandler(urgencyLevel UrgencyLevel, message string) {
	if l.options.LogLevel == UrgencyLevelNone {
		return
	}

	if urgencyLevel <= l.options.LogLevel {
		prefix := l.options.Prefix

		if l.templateInPrefix {
			prefix = time.Now().Local().Format(l.timeTemplate)
			if l.urgencyInTemplate {
				prefix = strings.Replace(prefix, "<{Q}>", l.urgencyList[urgencyLevel], -1)
			}
		}

		fmt.Fprintf(l.options.Destination, "%s%s\n", prefix, message)
	}
}

// Log - print message to log.
func (l *Daslog) Log(urgencyLevel UrgencyLevel, message string) {
	l.options.Handler(urgencyLevel, message)
}

// Logf - print message to log in printf style :)
func (l *Daslog) Logf(urgencyLevel UrgencyLevel, message string, a ...interface{}) {
	l.options.Handler(urgencyLevel, fmt.Sprintf(message, a...))
}

// Notice - print "notice" level message to log.
func (l *Daslog) Notice(message string) {
	l.options.Handler(UrgencyLevelNotice, message)
}

// Noticef - same as a Notice, but in printf style.
func (l *Daslog) Noticef(message string, a ...interface{}) {
	l.options.Handler(UrgencyLevelNotice, fmt.Sprintf(message, a...))
}

// Info - print "info" level message to log.
func (l *Daslog) Info(message string) {
	l.options.Handler(UrgencyLevelInfo, message)
}

// Infof - same as a Info, but in printf style.
func (l *Daslog) Infof(message string, a ...interface{}) {
	l.options.Handler(UrgencyLevelInfo, fmt.Sprintf(message, a...))
}

// Error - print "error" level message to log.
func (l *Daslog) Error(message string) {
	l.options.Handler(UrgencyLevelError, message)
}

// Errorf - same as a Error, but in printf style.
func (l *Daslog) Errorf(message string, a ...interface{}) {
	l.options.Handler(UrgencyLevelError, fmt.Sprintf(message, a...))
}

// Critical - print "critical" level message to log.
func (l *Daslog) Critical(message string) {
	l.options.Handler(UrgencyLevelCritical, message)
}

// Criticalf - same as a Critical, but in printf style.
func (l *Daslog) Criticalf(message string, a ...interface{}) {
	l.options.Handler(UrgencyLevelCritical, fmt.Sprintf(message, a...))
}
