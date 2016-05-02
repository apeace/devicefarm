package util

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
)

// CaptureWriter is an io.Writer that simply captures its inputs in a []string.
type CaptureWriter struct {
	out []string
}

// NewCaptureLogger returns a StandardLogger which will write its Print*() outputs
// to the returned CaptureWriter. This can be used to capture the outputs of Print*()
// methods instead of printing them to stdout.
func NewCaptureLogger() (*CaptureWriter, *StandardLogger) {
	capture := &CaptureWriter{}
	logger := NewStandardLogger(capture, os.Stderr)
	return capture, logger
}

func (w *CaptureWriter) Write(b []byte) (n int, err error) {
	n = len(b)
	w.out = append(w.out, string(b))
	return
}

func (w *CaptureWriter) Out() []string {
	return w.out[:]
}

// Logger is the interface used throughout devicefarm to write logs.
// The convention is that methods starting with Print write unformatted logs to sdout,
// while all other methods write formatted logs to stderr. Method signatures copied from Logrus:
// https://github.com/Sirupsen/logrus/blob/4b6ea7319e214d98c938f12692336f7ca9348d6b/logrus.go
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Warningln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
}

// StandardLogger is the basis for all loggers used in devicefarm.
type StandardLogger struct {
	out io.Writer
	Logger
}

// DefaultLogger writes to stdout and stderr
var DefaultLogger *StandardLogger = NewStandardLogger(os.Stdout, os.Stderr)

// NilLogger silences logs by writing to ioutil.Discard (equivalent of /dev/null)
var NilLogger *StandardLogger = NewStandardLogger(ioutil.Discard, ioutil.Discard)

// NewStandardLogger creates a StandardLogger with the given io.Writers. For Print*()
// methods, the out io.Writer is used. For non-Print methods, a standard logrus.Logger
// is used (and it is provided the err io.Writer).
func NewStandardLogger(out, err io.Writer) *StandardLogger {
	logrusLogger := logrus.New()
	logrusLogger.Out = err
	return &StandardLogger{out, logrusLogger}
}

func (logger *StandardLogger) Println(args ...interface{}) {
	str := fmt.Sprintln(args...)
	bytes := []byte(str)
	_, err := logger.out.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func (logger *StandardLogger) Printf(format string, args ...interface{}) {
	str := fmt.Sprintf(format, args...)
	bytes := []byte(str)
	_, err := logger.out.Write(bytes)
	if err != nil {
		panic(err)
	}
}

func (logger *StandardLogger) Print(args ...interface{}) {
	str := fmt.Sprint(args...)
	bytes := []byte(str)
	_, err := logger.out.Write(bytes)
	if err != nil {
		panic(err)
	}
}
