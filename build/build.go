/*

Package build provides data structures and functions to run app builds and
Device Farm tests.

*/
package build

import (
	"github.com/ride/devicefarm/config"
	"github.com/ride/devicefarm/util"
)

// A Build specifies all information needed to run a local app build: the
// working directory, the current Git branch of that directory, the full
// repo config, and the particular manifest for the given branch.
type Build struct {
	Dir      string
	Branch   string
	Config   *config.Config
	Manifest *config.BuildManifest
}

// Creates a new Build from a directory and a config file
func New(dir string, configFile string) (*Build, error) {
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
		Dir:      dir,
		Branch:   branch,
		Config:   config,
		Manifest: manifest,
	}
	return &build, nil
}

// Runs the build steps specified in this build's manifest, returning an error
// if any of the build steps produced an error
func (build *Build) RunLog(log util.Logger) error {
	_, err := util.RunAllLog(log, build.Dir, build.Manifest.Steps...)
	return err
}

func (build *Build) Run() error {
	return build.RunLog(util.NilLogger)
}
