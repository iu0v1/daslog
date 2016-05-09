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
		LogLevel:    UrgencyLevelCritical,
	})
	chkerr(err, t, "fail to init Daslog")

	l.Notice("Notice")
	l.Noticef("Noticef")
	l.Info("Info")
	l.Infof("Infof")
	l.Error("Error")
	l.Errorf("Errorf")
	l.Critical("Critical")
	l.Criticalf("Criticalf")

	outTests := []outTest{
		outTest{N: "Notice", E: "Notice"},
		outTest{N: "Noticef", E: "Noticef"},
		outTest{N: "Info", E: "Info"},
		outTest{N: "Infof", E: "Infof"},
		outTest{N: "Error", E: "Error"},
		outTest{N: "Errorf", E: "Errorf"},
		outTest{N: "Critical", E: "Critical"},
		outTest{N: "Criticalf", E: "Criticalf"},
	}

	outTestsBuf := strings.Split(buf.String(), "\n")

	for i, tst := range outTestsBuf {
		if i == len(outTestsBuf)-1 {
			continue
		}

		if tst != outTests[i].E {
			t.Fatalf("'%s' != '%s' (out test: #%d, name: %s)\n", outTests[i], tst, i, outTests[i].N)
		}
	}

	prfxTests := []prfxTmplTest{
		prfxTmplTest{
			E: time.Now().Local().Format("2006-01-02"),
			O: Options{Prefix: "{{.F}} ", LogLevel: UrgencyLevelCritical},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("2006"),
			O: Options{Prefix: "{{.Y}} ", LogLevel: UrgencyLevelCritical},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("06"),
			O: Options{Prefix: "{{.y}} ", LogLevel: UrgencyLevelCritical},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("01"),
			O: Options{Prefix: "{{.m}} ", LogLevel: UrgencyLevelCritical},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("Jan"),
			O: Options{Prefix: "{{.b}} ", LogLevel: UrgencyLevelCritical},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("January"),
			O: Options{Prefix: "{{.B}} ", LogLevel: UrgencyLevelCritical},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("02"),
			O: Options{Prefix: "{{.d}} ", LogLevel: UrgencyLevelCritical},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("Mon"),
			O: Options{Prefix: "{{.a}} ", LogLevel: UrgencyLevelCritical},
		},
		prfxTmplTest{
			E: time.Now().Local().Format("PM"),
			O: Options{Prefix: "{{.p}} ", LogLevel: UrgencyLevelCritical},
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
		LogLevel:     UrgencyLevelCritical,
	})
	chkerr(err, t, "fail to init Daslog")

	l3.Info("test")

	if buf.String() != "test\n" || buf2.String() != "test\n" {
		t.Fatalf("destinations test fail\n")
	}
}
