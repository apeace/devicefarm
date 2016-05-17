/*

Package build provides data structures and functions to run app builds and
Device Farm tests.

*/
package build

import (
	"github.com/ride/devicefarm/awsutil"
	"github.com/ride/devicefarm/config"
	"github.com/ride/devicefarm/util"
)

// A Build specifies all information needed to run a local app build: the
// working directory, the current Git branch of that directory, the full
// repo config, and the particular manifest for the given branch.
type Build struct {
	Log      util.Logger
	Dir      string
	Branch   string
	Config   *config.Config
	Manifest *config.BuildManifest
	Client   *awsutil.DeviceFarm
}

// Creates a new Build from a directory and a config file
func New(log util.Logger, dir string, configFile string) (*Build, error) {
	config, err := config.New(configFile)
	if err != nil {
		return nil, err
	}
	branch, err := util.GitBranch(dir)
	if err != nil {
		return nil, err
	}
	manifest := config.BranchManifest(branch)
	if runnable, err := manifest.IsRunnable(); !runnable {
		return nil, err
	}
	build := Build{
		Log:      log,
		Dir:      dir,
		Branch:   branch,
		Config:   config,
		Manifest: manifest,
	}
	return &build, nil
}

func (build *Build) Run() error {
	// Runs the build steps specified in this build's manifest, returning an error
	// if any of the build steps produced an error
	_, err := util.RunAllLog(build.Log, build.Dir, build.Manifest.Steps...)
	return err
}
