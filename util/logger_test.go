package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCaptureLogger(t *testing.T) {
	assert := assert.New(t)
	out, log := NewCaptureLogger()
	log.Println("foo", "bar")
	log.Printf("baz %s", "buzz")
	log.Print("blah", "boo")
	expected := []string{
		"foo bar\n",
		"baz buzz",
		"blahboo",
	}
	assert.Equal(expected, out.Out())
}

type errorWriter struct{}

func (w *errorWriter) Write(b []byte) (n int, err error) {
	return 0, fmt.Errorf("Fake error")
}

func TestStandardLogger(t *testing.T) {
	assert := assert.New(t)
	out := &errorWriter{}
	log := NewStandardLogger(out, out)
	println := func() { log.Println("foo") }
	printf := func() { log.Printf("foo") }
	print := func() { log.Print("foo") }
	for _, f := range []func(){println, printf, print} {
		assert.Panics(f)
	}
}

func ExampleCaptureWriter() {
	w := &CaptureWriter{}
	w.Write([]byte("foo"))
	w.Write([]byte("bar"))
	w.Out() // []string{"foo", "bar"}
}

func ExampleNewCaptureLogger() {
	w, log := NewCaptureLogger()
	log.Println("foo") // nothing written to stdout
	log.Debugln("bar") // "bar" is written to stderr
	log.Println("baz") // nothing written to stdout
	w.Out()            // []string{"foo", "baz"}
}

func ExampleNewStandardLogger() {
	log := NewStandardLogger(os.Stdout, os.Stderr)
	log.Println("foo") // unformatted log to stdout
	log.Debugln("bar") // logrus formatted log to stderr
}
