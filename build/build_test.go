package build

import (
	"github.com/ride/devicefarm/config"
	"github.com/ride/devicefarm/util"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

var log *util.StandardLogger = util.NilLogger

func TestNew(t *testing.T) {
	assert := assert.New(t)

	tmpDir, err := ioutil.TempDir("", "devicefarm")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)

	configFile := "devicefarm.yml"
	absConfigFile := path.Join(tmpDir, configFile)

	// at this point we should fail because the dir has no config file
	build, err := New(log, tmpDir, absConfigFile)
	assert.Nil(build)
	assert.NotNil(err)

	// this config file will fail for the "foobar" branch because it
	// has no default devicepool. but it will succeed for the "master" branch
	util.CopyFile("../config/testdata/config_nofoobar.yml", absConfigFile)

	// at this point we should fail because the dir is not a git repo
	build, err = New(log, tmpDir, absConfigFile)
	assert.Nil(build)
	assert.NotNil(err)

	util.RunAll(tmpDir,
		"git init",
		"git config user.email 'devops@ride.com'",
		"git config user.name 'Devops'",
		"git checkout -b foobar",
		"git add "+configFile,
		"git commit "+configFile+" -m foo")

	// at this point we should fail because the "foobar" manifest is
	// not runnable
	build, err = New(log, tmpDir, absConfigFile)
	assert.Nil(build)
	assert.NotNil(err)

	util.RunAll(tmpDir, "git checkout -b master")

	// now we should succeed
	build, err = New(log, tmpDir, absConfigFile)
	assert.Nil(err)
	assert.NotNil(build)
}

func TestBuildRun(t *testing.T) {
	assert := assert.New(t)

	tmpDir, err := ioutil.TempDir("", "devicefarm")
	if err != nil {
		t.Error(err)
	}
	defer os.RemoveAll(tmpDir)

	// should succeed. we don't need a complete Build in order to Run()
	build := Build{
		Log: log,
		Dir: tmpDir,
		Manifest: &config.BuildManifest{
			Steps:      []string{"echo Foo", "echo Bar"},
			DevicePool: "",
		},
	}
	err = build.Run()
	assert.Nil(err)

	// should fail because "exit 1" produces an error
	build.Manifest.Steps = []string{"echo Foo", "exit 1"}
	err = build.Run()
	assert.NotNil(err)
}
