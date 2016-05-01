package util

import (
	"errors"
	"github.com/stretchr/testify/assert"
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
	return 0, errors.New("Fake error")
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
