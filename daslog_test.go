package daslog

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func chkerr(e error, t *testing.T, m string) {
	if e != nil {
		_, fn, line, _ := runtime.Caller(1)
		fnb := filepath.Base(fn)
		t.Fatalf("%s\ncall: %s:%d\nerr : %v\n", m, fnb, line, e)
	}
}

type prfxTmplTest struct {
	E string  // expected
	O Options // options
}

type outTest struct {
	N string // name
	E string // expected
}

// TODO : add more tests
func TestMain(t *testing.T) {
	var buf bytes.Buffer

	l, err := New(Options{
		Destination: &buf,
		LogLevel:    UrgencyLevelNotice,
	})
	chkerr(err, t, "fail to init Daslog")

	l.Notice("Notice")
	l.Noticef("Noticef\n")
	l.Info("Info")
	l.Infof("Infof\n")
	l.Error("Error")
	l.Errorf("Errorf\n")
	l.Critical("Critical")
	l.Criticalf("Criticalf\n")

	outTests := []outTest{
		outTest{N: "Notice", E: "Notice\n"},
		outTest{N: "Noticef", E: "Noticef\n"},
		outTest{N: "Info", E: "Info\n"},
		outTest{N: "Infof", E: "Infof\n"},
		outTest{N: "Error", E: "Error\n"},
		outTest{N: "Errorf", E: "Errorf\n"},
		outTest{N: "Critical", E: "Critical\n"},
		outTest{N: "Criticalf", E: "Criticalf\n"},
	}

	outTestsBuf := strings.SplitAfterN(buf.String(), "\n", 8)

	for i, tst := range outTestsBuf {
		if tst != outTests[i].E {
			t.Fatalf("'%s' != '%s' (out test: #%d, name: %s)\n", outTests[i].E, tst, i, outTests[i].N)
		}
	}

	prfxTests := []prfxTmplTest{
		prfxTmplTest{
			E: time.Now().Local().Format("2006-01-02"),
			O: Options{Prefix: "{{.F}} ", LogLevel: UrgencyLevelNotice},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("2006"),
			O: Options{Prefix: "{{.Y}} ", LogLevel: UrgencyLevelNotice},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("06"),
			O: Options{Prefix: "{{.y}} ", LogLevel: UrgencyLevelNotice},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("01"),
			O: Options{Prefix: "{{.m}} ", LogLevel: UrgencyLevelNotice},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("Jan"),
			O: Options{Prefix: "{{.b}} ", LogLevel: UrgencyLevelNotice},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("January"),
			O: Options{Prefix: "{{.B}} ", LogLevel: UrgencyLevelNotice},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("02"),
			O: Options{Prefix: "{{.d}} ", LogLevel: UrgencyLevelNotice},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("Mon"),
			O: Options{Prefix: "{{.a}} ", LogLevel: UrgencyLevelNotice},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("PM"),
			O: Options{Prefix: "{{.p}} ", LogLevel: UrgencyLevelNotice},
		},
	}

	for i, tst := range prfxTests {
		buf.Reset()
		tst.O.Destination = &buf

		l2, e := New(tst.O)
		chkerr(e, t, fmt.Sprintf("fail to init Daslog (date test #%d)", i))

		l2.Info("test")

		out := strings.Replace(buf.String(), "\n", "", -1)
		exp := tst.E + " test"

		if out != exp {
			t.Fatalf("'%s' != '%s' (date test #%d)\n", out, exp, i)
		}
	}

	buf.Reset()
	var buf2 bytes.Buffer

	l3, err := New(Options{
		Destinations: []io.Writer{&buf, &buf2},
		LogLevel:     UrgencyLevelNotice,
	})
	chkerr(err, t, "fail to init Daslog")

	l3.Info("test")

	if buf.String() != "test\n" || buf2.String() != "test\n" {
		t.Fatalf("destinations test fail; buf == '%s', buf2 == '%s'\n",
			buf.String(), buf2.String(),
		)
	}

	// UL tests
	buf.Reset()
	l4, err := New(Options{
		Destinations: []io.Writer{&buf},
		LogLevel:     UrgencyLevelNone,
	})
	chkerr(err, t, "fail to init Daslog")

	l4.Critical("test")
	if buf.String() != "" {
		t.Fatalf("non empty buf: %s \n", buf.String())
	}

	buf.Reset()
	l4, err = New(Options{
		Destinations: []io.Writer{&buf},
		LogLevel:     UrgencyLevelError,
	})
	chkerr(err, t, "fail to init Daslog")

	l4.Notice("notice")
	l4.Info("info")
	l4.Error("error")
	l4.Critical("critical")

	outTestsBuf = strings.SplitAfterN(buf.String(), "\n", 2)

	if len(outTestsBuf) != 2 {
		t.Fatalf("UL test error #0: %d != 2 \n", len(outTestsBuf))
	}
}
