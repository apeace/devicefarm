package build

import (
	"github.com/ride/devicefarm/config"
	"github.com/ride/devicefarm/util"
	"log"
)

type Build struct {
	Dir      string
	Branch   string
	Config   *config.Config
	Manifest *config.BuildManifest
}

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

func (build *Build) Run() error {
	outputs := util.RunAll(build.Dir, []string(build.Manifest.Steps)...)
	for _, output := range outputs {
		log.Println("$ " + output.Cmd)
		log.Println(output.Output)
		if output.Err != nil {
			log.Println(output.Err)
			return output.Err
		}
	}
	return nil
}
