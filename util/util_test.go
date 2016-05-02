package util

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestGitBranch(t *testing.T) {
	assert := assert.New(t)
	tmpDir, err := ioutil.TempDir("", "devicefarm")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)
	// at this point we should have an error because
	// tmpDir is not a git repository
	branch, err := GitBranch(tmpDir)
	assert.NotNil(err)
	// create git repository, a branch, and commit a file
	RunAll(tmpDir,
		"git init",
		"git config user.email 'devops@ride.com'",
		"git config user.name 'Devops'",
		"git checkout -b foobar",
		"touch foo",
		"git add foo",
		"git commit foo -m foo")
	branch, err = GitBranch(tmpDir)
	assert.Nil(err)
	assert.Equal("foobar", branch)
	// get into "detached" state by adding another commit
	// and then checking out HEAD^ (meaning previous commit)
	RunAll(tmpDir,
		"touch bar",
		"git add bar",
		"git commit bar -m bar",
		"git checkout HEAD^")
	branch, err = GitBranch(tmpDir)
	assert.Equal(ErrGitDetached, err)
}

func TestCmd(t *testing.T) {
	assert := assert.New(t)
	cmd := Cmd("/dir", "echo bar baz")
	assert.Equal("echo", path.Base(cmd.Path))
	assert.Equal([]string{"echo", "bar", "baz"}, cmd.Args)
	assert.Equal("/dir", cmd.Dir)
}

func TestRunAll(t *testing.T) {
	assert := assert.New(t)
	tmpDir, err := ioutil.TempDir("", "devicefarm")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)
	out, log := NewCaptureLogger()
	outputs, err := RunAllLog(log, tmpDir,
		"echo Foo",
		"exit 1",
		"echo Bar")
	assert.NotNil(err)
	assert.Equal(2, len(outputs))
	assert.Equal(CmdOutput{"echo Foo", "Foo", nil}, *outputs[0])
	assert.NotNil(outputs[1].Err)
	assert.Equal([]string{"$ echo Foo\n", "$ exit 1\n"}, out.Out())
}

func TestCopyFile(t *testing.T) {
	assert := assert.New(t)
	tmpDir, err := ioutil.TempDir("", "devicefarm")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)

	dstFilename := path.Join(tmpDir, "foo.txt")

	// should fail because source file does not exist
	err = CopyFile("testdata/does-not-exist.txt", dstFilename)
	assert.NotNil(err)

	// should fail because blank filename cannot be created
	err = CopyFile("testdata/foo.txt", "")
	assert.NotNil(err)

	// should succeed
	err = CopyFile("testdata/foo.txt", dstFilename)
	assert.Nil(err)
	bytes, err := ioutil.ReadFile(dstFilename)
	assert.Nil(err)
	assert.Equal("Foo\n", string(bytes))
}
